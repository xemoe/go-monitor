package monitor

import log "github.com/sirupsen/logrus"

func (monitor *Monitor) Println(message string) {
	if monitor.writeToConsole {
		log.Println(message)
	}
}

func (monitor *Monitor) Printf(message string, a ...interface{}) {
	if monitor.writeToConsole {
		log.Printf(message, a...)
	}
}

func (monitor *Monitor) Errorf(message string, a ...interface{}) {
	if monitor.writeToConsole {
		log.Errorf(message, a...)
	}
}

func (monitor *Monitor) Warnf(message string, a ...interface{}) {
	if monitor.writeToConsole {
		log.Warnf(message, a...)
	}
}

func (monitor *Monitor) Debugf(message string, a ...interface{}) {
	if monitor.writeToConsole {
		log.Debugf(message, a...)
	}
}
