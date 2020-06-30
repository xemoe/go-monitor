package gomonitor


import (

	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
    "github.com/rob121/go-monitor/ndriver"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/mitchellh/go-ps"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)


var DriverRegistry map[string]interface{}


type NotificationDriver interface {
    Init() error
    Alert(proc string, server ServerInfo, state string, c *cache.Cache) error
}

type ServerInfo struct{ 
    Host string
    Ip string
}

func (s ServerInfo) String() string{
    
    return  s.Host + " with IP " + s.Ip
    
}

// All these values should be in lowercase in the yaml file
// See sample.go-monitor.yaml
type Monitor struct {
	Processes []string
	Config struct {
        DefaultTTLSeconds     time.Duration
		NotifyServiceReturn   bool
		CheckFrequencySeconds time.Duration
	}
    NotificationDriver string
    Drivers map[string]interface{}
	writeToConsole bool
}

//add driver to this map
func init(){
    
    DriverRegistry = make(map[string]interface{})
    
    DriverRegistry["messagebird"]=MessagebirdDriver{}
    DriverRegistry["pushover"]=PushoverDriver{}
    
}

func Start(configPath string) {

   
	monitor, err := createMonitorFromFile(configPath)
	//set default 
	monitor.writeToConsole = true

	monitor.Println("Go Monitor running")
  
    if val,ok := DriverRegistry[monitor.NotificationDriver]; ok{

      
       err := ndriver.Call(val,"Init")
       
       if(err!=nil){log.Println(err)}  
    
    }else{
           
         log.Fatal("Driver  Not Available")  
           
    }
    
	// Default notification from config
	// Refresh time is 60 seconds
	c := cache.New(monitor.Config.DefaultTTLSeconds*time.Second, 60*time.Second)
	procErrChan := make(chan string, len(monitor.Processes))
	procSuccessChan := make(chan string, len(monitor.Processes))

	server, err := monitor.getServerInfo()
	if err != nil {
		monitor.Println("Error getting server information, using NIL")
		server =  ServerInfo{}
	}

	// Parent waitgroup for the two go functions below
	var wgParent sync.WaitGroup
	wgParent.Add(3)

	// One go func for the adding of procs to error channel
	go func() {
		var wg sync.WaitGroup
		for {
			wg.Add(len(monitor.Processes))

			for index, proc := range monitor.Processes {
				go monitor.checkProcess(proc, procErrChan,procSuccessChan, &wg)

				// Sleep when the loop is done
				// This is how often the checks for each process will run
				if index == len(monitor.Processes)-1 {
					// Check every 60 seconds
					time.Sleep(monitor.Config.CheckFrequencySeconds * time.Second)
				}
			}
			wg.Wait()
		}
	}()

	// Another go func for reading the results from the error chan
	go func() {
		for {
			procErr := <-procErrChan
			go monitor.notifyProcError(procErr, server, c)
		}
	}()
	
	
    go func() {
		for {
			procSuccess := <-procSuccessChan
			go monitor.notifyProcSuccess(procSuccess, server, c)
		}
	}()

	// Never die
	wgParent.Wait()
}


func load(path string) *viper.Viper {
	v := viper.New()
	v.SetConfigName("config") // name of config file (without extension)

	if len(path) > 0 {
		v.AddConfigPath(path)
	}

	v.AddConfigPath("/etc/go-monitor/") // path to look for the config file in
	v.AddConfigPath(".") // optionally look for config in the working directory
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)

		if len(v.AllSettings()) == 0 {
			fmt.Println("*** Invalid Config, you have an error ****")
		}
	})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(err)
		}
	}

	return v
}

func (monitor *Monitor) Println(message string) {

	if monitor.writeToConsole {
		log.Println(message)
	}
}

func (monitor *Monitor) Printf(message string, a ...interface{}) {

	if monitor.writeToConsole {
		log.Printf(message, a)
	}
}

func createMonitorFromFile(configFile string) (monitor *Monitor, err error) {

  

    config := load(configFile)

    err = config.Unmarshal(&monitor)
  
    if err != nil {
        panic(err)
    }

	err = monitor.validate()

	return monitor,err
}

func (monitor *Monitor) validate() error {
	// Do validation checks
	if len(monitor.Processes) < 1 {
		return errors.New("Config: We need to monitor at least one process")
	} else {
		monitor.Printf("Processes %s\n", monitor.Processes)
	}

	if monitor.Config.DefaultTTLSeconds == 0 {
		monitor.Config.DefaultTTLSeconds = 30000
	}
	if monitor.Config.CheckFrequencySeconds == 0 {
		monitor.Config.CheckFrequencySeconds = 60
	}


	monitor.Printf("DefaultTTLSeconds %d\n", monitor.Config.DefaultTTLSeconds)
	monitor.Printf("CheckFrequencySeconds %d\n", monitor.Config.CheckFrequencySeconds)

	return nil
}

func (monitor *Monitor) getServerInfo() (server ServerInfo, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ServerInfo{}, err
	}

	var ip net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ServerInfo{}, err
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
		}
	}

	// Get hostname
	host, err := os.Hostname()
	if err != nil {
		return ServerInfo{}, err
	}

	return ServerInfo{Host: host,Ip: ip.String()},nil

}

func (monitor *Monitor) checkProcess(processName string, procErrChan chan string,procSuccessChan chan string, wg *sync.WaitGroup) {


 
	if strings.HasPrefix(processName, "tcp://") {
		monitor.checkTcpSocket(strings.TrimPrefix(processName, "tcp://"), procErrChan, procSuccessChan, wg)
	} else if strings.HasPrefix(processName, "http://") || strings.HasPrefix(processName, "https://") {
		monitor.checkHttpEndpoint(processName, procErrChan,procSuccessChan, wg)
	} else {
		monitor.checkLocalProcess(processName, procErrChan,procSuccessChan, wg)
	}
}

func (monitor *Monitor) checkTcpSocket(tcpAddress string, procErrChan chan string, procSuccessChan chan string, wg *sync.WaitGroup) {
	monitor.Printf("Checking for tcp socket %s\n", tcpAddress)

	conn, err := net.Dial("tcp", tcpAddress)
	defer wg.Done()
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	if err != nil {
		monitor.Printf("Error: unable to open socket! %s\n", tcpAddress)
		procErrChan <- tcpAddress
        procSuccessChan <- ""
        return
	} else {
		monitor.Printf("Successful connection to %s \n", tcpAddress)
	}

	// Doing this keeps the channel open
	// If this is not done, the channel closes and there is a fatal error
	procErrChan <- ""
	procSuccessChan <- tcpAddress
}

func (monitor *Monitor) checkHttpEndpoint(httpEndpoint string, procErrChan chan string, procSuccessChan chan string, wg *sync.WaitGroup) {
	monitor.Printf("Checking http endpoint %s\n", httpEndpoint)
    defer wg.Done()
	resp, err := http.DefaultClient.Get(httpEndpoint)

	if err != nil {
		monitor.Printf("Error: unable to connect to %s - %s\n", httpEndpoint, err.Error())
		procErrChan <- httpEndpoint
        procSuccessChan <- ""
        return
	} else if resp.Status != "200 OK" {
		monitor.Printf("Error: non 200 status from %s - %s\n", httpEndpoint, resp.Status)
		procErrChan <- httpEndpoint
        procSuccessChan <- ""
        return
	} else {
		monitor.Printf("%s returns 200 OK\n", httpEndpoint, resp.Status)
	}

	// Doing this keeps the channel open
	// If this is not done, the channel closes and there is a fatal error
	procErrChan <- ""
	procSuccessChan <- httpEndpoint
}

func (monitor *Monitor) checkLocalProcess(processName string, procErrChan chan string, procSuccessChan chan string, wg *sync.WaitGroup) {
	monitor.Printf("Checking for process %s\n", processName)
    
    defer wg.Done()
    
    pid,_,_ := findProcess(processName)

	if pid == 0  {
		monitor.Printf("Error: no process %s found running!\n", processName)
		procErrChan <- processName
		procSuccessChan <- ""
		return
	}

	// Doing this keeps the channel open
	// If this is not done, the channel closes and there is a fatal error
	procErrChan <- ""
	procSuccessChan <- processName
}



func findProcess(key string) (int, string, error) {
    pname := ""
    pid := 0
    
    err := errors.New("not found")
    
    ps, _ := ps.Processes()
    
    for i, _ := range ps {
        
 
        
        if ps[i].Executable() == key {
            pid = ps[i].Pid()
            pname = ps[i].Executable()
            err = nil
            break
        }
    }
    
   
    
    return pid, pname, err
} 

// Notifyproceerror sends a notification for a given process
func (monitor *Monitor) notifyProcError(proc string, server ServerInfo, c *cache.Cache) {
	if len(proc) > 0 {
		monitor.Printf("### ERROR: proc %s not running!\n", proc)

		// Check cache for process
		_, found := c.Get(proc)
		if found {
			// Wait until expiry before another notification
			monitor.Printf("Alert previously sent for  %s, skipping...\n", proc)
			return
		}

		// If proc not in cache, store in cache
		c.Set(proc, true, cache.DefaultExpiration)

		// Send text message

       //plugin stuff here

       if val,ok := DriverRegistry[monitor.NotificationDriver]; ok{
          
       err := ndriver.Call(val,"Alert",monitor,proc,server,"DOWN")
       
       if(err!=nil){
           log.Printf("Error: %s",err)
           return
           
       }  
    
       }else{
           
         log.Println("Driver Not Available")  
         
         return
           
       }
		
		
		monitor.Println("Down Notification sent!")
	}
}

func (monitor *Monitor) notifyProcSuccess(proc string, server ServerInfo, c *cache.Cache) {
    
    if(!monitor.Config.NotifyServiceReturn){
        //disabled notifications when the service comes back "up"
        return
    }
    

	if len(proc) > 0 {
		monitor.Printf("### Success: proc %s is running!\n", proc)

		// Check cache for process
		_, found := c.Get(proc)
		if !found {
			// Wait until expiry before another notification
			monitor.Printf("Process %s up, no notification for up sent\n", proc)
			return
		}

		// Send text message

       //plugin stuff here

       if val,ok := DriverRegistry[monitor.NotificationDriver]; ok{
          
       err := ndriver.Call(val,"Alert",monitor,proc,server,"UP")
       
       if(err!=nil){
           log.Printf("Error: %s",err)
           return
           
       }  
    
       }else{
           
         log.Println("Driver Not Available")  
         
         return
           
       }
		
       // If proc  in cache, delete, indicator its 'UP'
		c.Delete(proc)
		monitor.Println("Up Notification sent!")
	}
}
