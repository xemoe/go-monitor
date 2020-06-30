package main

import (
    "flag"
	//"fmt"
	"github.com/rob121/go-monitor/monitor"
)


var configFile string

func main(){
    
    
    //defaultConfigfile := "/usr/local/etc/go-monitor.yml"
	flag.StringVar(&configFile,"f","", "config file path")
///	writeToConsole := flag.Bool("o", false, fmt.Sprintf("output, if true will write to console"))
	//flag.Parse()
    
    gomonitor.Start(configFile)
    
    
}