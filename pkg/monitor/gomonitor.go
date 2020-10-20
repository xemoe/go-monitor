package monitor

import (
	"errors"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-ps"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/xemoe/go-monitor/pkg/ndriver"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

var DriverRegistry map[string]interface{}

type r interface {
	Init() error
	Alert(proc string, server ServerInfo, state string, c *cache.Cache) error
}

//
// All these values should be in lowercase in the yaml file
// See config.yml.default
//
type Monitor struct {
	Processes []string
	Config    struct {
		DefaultTTLSeconds     time.Duration
		NotifyServiceReturn   bool
		CheckFrequencySeconds time.Duration
	}
	NotificationDriver string
	Drivers            map[string]interface{}
	writeToConsole     bool
}

//
// Add Notification Driver to this map
// เพิ่ม Notification Driver
//
func init() {

	log.SetFormatter(&nested.Formatter{
		HideKeys:       false,
		NoColors:       false,
		NoFieldsColors: false,
		ShowFullLevel:  true,
		FieldsOrder:    []string{"process", "status"},
	})

	log.SetLevel(log.DebugLevel)

	//
	// Enable to show caller
	//
	// log.SetReportCaller(true)
}

//
// SetDriverRegistry from main
// เพิ่ม Notification Driver
//
func SetDriverRegistry(r map[string]interface{}) {
	DriverRegistry = r
}

//
// Start monitoring
//
func Start(configPath string) {

	monitor, err := createMonitorFromFile(configPath)

	//
	// Set default
	// กำหนดค่าเริ่มต้นของ Monitor
	//
	monitor.writeToConsole = true

	log.WithFields(log.Fields{
		"status": STATUS_START,
	}).Info("Go Monitor running")

	if val, ok := DriverRegistry[monitor.NotificationDriver]; ok {

		err := ndriver.Call(val, "Init")
		if err != nil {
			log.Println(err)
		}

	} else {
		log.Fatal("Driver  Not Available")
	}

	//
	// Default notification from config
	// Refresh time is 60 seconds
	// สร้าง channel ตามจำนวน process ที่ต้องการ monitor
	//
	c := cache.New(monitor.Config.DefaultTTLSeconds*time.Second, 60*time.Second)
	procErrChan := make(chan string, len(monitor.Processes))
	procSuccessChan := make(chan string, len(monitor.Processes))

	server, err := monitor.getServerInfo()
	if err != nil {
		log.Errorf("Error: getting server information, using NIL")
		server = ServerInfo{}
	}

	//
	// Parent waitgroup for the two go functions below
	//
	var wgParent sync.WaitGroup
	wgParent.Add(3)

	//
	// One go func for the adding of procs to error channel
	//
	go func() {
		var wg sync.WaitGroup
		for {
			wg.Add(len(monitor.Processes))

			for index, proc := range monitor.Processes {
				go monitor.checkProcess(proc, procErrChan, procSuccessChan, &wg)

				//
				// Sleep when the loop is done
				// This is how often the checks for each process will run
				//
				if index == len(monitor.Processes)-1 {
					//
					// Check every 60 seconds
					//
					time.Sleep(monitor.Config.CheckFrequencySeconds * time.Second)
				}
			}
			wg.Wait()
		}
	}()

	//
	// Another go func for reading the results from the error chan
	//
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

	//
	// Never die
	//
	wgParent.Wait()
}

func load(path string) *viper.Viper {
	v := viper.New()

	//
	// name of config file without extension
	// E.g. "config" for "config.yml"
	//
	v.SetConfigName("config")

	if len(path) > 0 {
		v.AddConfigPath(path)
	}

	//
	// path to look for the config file in
	// ค่าพื้นฐานของตำแหน่งที่เก็บ configuration file
	//
	v.AddConfigPath("/etc/go-monitor/")

	//
	// optionally look for config in the working directory
	// หรืออ่าน configuration file จากตำแหน่งปัจจุบัน
	//
	v.AddConfigPath(".")
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {

		log.Warnf("Config file changed: %s", e.Name)

		if len(v.AllSettings()) == 0 {
			log.Error("*** Invalid Config, you have an error ****")
		}
	})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//
			// Config file not found; ignore error if desired
			//
			panic(err)
		}
	}

	return v
}

func createMonitorFromFile(configFile string) (monitor *Monitor, err error) {

	config := load(configFile)

	err = config.Unmarshal(&monitor)
	if err != nil {
		panic(err)
	}

	err = monitor.Validate()

	return monitor, err
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
