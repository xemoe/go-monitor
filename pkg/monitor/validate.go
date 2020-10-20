package monitor

import "errors"

func (monitor *Monitor) Validate() error {
	//
	// Do validation checks
	//
	if len(monitor.Processes) < 1 {
		return errors.New("Config: We need to monitor at least one process")
	} else {
		monitor.Printf("Processes %s", monitor.Processes)
	}

	if monitor.Config.DefaultTTLSeconds == 0 {
		monitor.Config.DefaultTTLSeconds = 30000
	}
	if monitor.Config.CheckFrequencySeconds == 0 {
		monitor.Config.CheckFrequencySeconds = 60
	}

	monitor.Printf("DefaultTTLSeconds %d", monitor.Config.DefaultTTLSeconds)
	monitor.Printf("CheckFrequencySeconds %d", monitor.Config.CheckFrequencySeconds)

	return nil
}
