package monitor

import (
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/xemoe/go-monitor/pkg/ndriver"
)

//
// NotifyProcError sends a notification for a given process
//
func (monitor *Monitor) notifyProcError(proc string, server ServerInfo, c *cache.Cache) {

	if len(proc) > 0 {

		log.WithFields(log.Fields{
			"process": proc,
			"status":  STATUS_NOTIFY_DOWN,
		}).Warnf("Error, process=%s not running.", proc)

		//
		// Check cache for process
		//
		_, found := c.Get(proc)
		if found {
			//
			// Wait until expiry before another notification
			//
			log.WithFields(log.Fields{
				"process": proc,
				"status":  STATUS_NOTIFY_SKIP,
			}).Warnf("Alert, process=%s previously sent skipping.", proc)

			return
		}

		//
		// If proc not in cache, store in cache
		//
		c.Set(proc, true, cache.DefaultExpiration)

		//
		// Send text message
		//
		// Plugin stuff here
		//
		if val, ok := DriverRegistry[monitor.NotificationDriver]; ok {

			err := ndriver.Call(val, "Alert", monitor, proc, server, "DOWN")
			if err != nil {
				log.Errorf("Error: %s", err)
				return
			}

		} else {
			log.Errorf("Driver Not Available")
			return
		}

		log.WithFields(log.Fields{
			"process": proc,
			"status":  STATUS_NOTIFY_DOWN_SENT,
		}).Infof("Down, process=%s notification sent.", proc)
	}
}

func (monitor *Monitor) notifyProcSuccess(proc string, server ServerInfo, c *cache.Cache) {

	if !monitor.Config.NotifyServiceReturn {
		//
		// Disabled notifications when the service comes back "up"
		//
		return
	}

	if len(proc) > 0 {

		log.WithFields(log.Fields{
			"process": proc,
			"status":  STATUS_NOTIFY_RUNNING,
		}).Infof("Success, process=%s is running.", proc)

		//
		// Check cache for process
		//
		_, found := c.Get(proc)
		if !found {
			//
			// Wait until expiry before another notification
			//
			log.WithFields(log.Fields{
				"process": proc,
				"status":  STATUS_NOTIFY_UP,
			}).Infof("Up, process=%s no notification for up sent.", proc)

			return
		}

		//
		// Send text message
		//
		// Plugin stuff here
		//
		if val, ok := DriverRegistry[monitor.NotificationDriver]; ok {

			err := ndriver.Call(val, "Alert", monitor, proc, server, "UP")
			if err != nil {
				log.Errorf("Error: %s", err)
				return
			}

		} else {
			log.Errorf("Driver Not Available")
			return
		}

		log.WithFields(log.Fields{
			"process": proc,
			"status":  STATUS_NOTIFY_UP_SENT,
		}).Infof("Up, process=%s notification sent", proc)

		//
		// If proc  in cache, delete, indicator its 'UP'
		//
		c.Delete(proc)
	}
}
