package main

import (
	"flag"

	"github.com/xemoe/go-monitor/pkg/drivers"
	"github.com/xemoe/go-monitor/pkg/monitor"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "f", "", "config file path")

	monitor.SetDriverRegistry(setDriverRegistry())
	monitor.Start(configFile)
}

func setDriverRegistry() map[string]interface{} {

	registry := make(map[string]interface{})

	registry["messagebird"] = drivers.MessagebirdDriver{}
	registry["pushover"] = drivers.PushoverDriver{}
	registry["webhook"] = drivers.WebHookDriver{}

	return registry
}
