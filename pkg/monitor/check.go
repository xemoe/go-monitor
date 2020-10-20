package monitor

import (
	"net"
	"net/http"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func (monitor *Monitor) checkProcess(processName string, procErrChan chan string, procSuccessChan chan string, wg *sync.WaitGroup) {
	if strings.HasPrefix(processName, "tcp://") {
		monitor.checkTcpSocket(strings.TrimPrefix(processName, "tcp://"), procErrChan, procSuccessChan, wg)
	} else if strings.HasPrefix(processName, "http://") || strings.HasPrefix(processName, "https://") {
		monitor.checkHttpEndpoint(processName, procErrChan, procSuccessChan, wg)
	} else {
		monitor.checkLocalProcess(processName, procErrChan, procSuccessChan, wg)
	}
}

func (monitor *Monitor) checkTcpSocket(tcpAddress string, procErrChan chan string, procSuccessChan chan string, wg *sync.WaitGroup) {

	monitor.Printf("Checking for tcp socket %s", tcpAddress)

	conn, err := net.Dial("tcp", tcpAddress)
	defer wg.Done()
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	if err != nil {
		monitor.Errorf("Error: unable to open socket! %s", tcpAddress)
		procErrChan <- tcpAddress
		procSuccessChan <- ""
		return
	} else {
		monitor.Printf("Successful connection to %s", tcpAddress)
	}

	//
	// Doing this keeps the channel open
	// If this is not done, the channel closes and there is a fatal error
	//
	procErrChan <- ""
	procSuccessChan <- tcpAddress
}

func (monitor *Monitor) checkHttpEndpoint(httpEndpoint string, procErrChan chan string, procSuccessChan chan string, wg *sync.WaitGroup) {

	monitor.Printf("Checking http endpoint %s", httpEndpoint)
	defer wg.Done()
	resp, err := http.DefaultClient.Get(httpEndpoint)

	if err != nil {
		monitor.Errorf("Error: unable to connect to %s - %s", httpEndpoint, err.Error())
		procErrChan <- httpEndpoint
		procSuccessChan <- ""
		return
	} else if resp.Status != "200 OK" {
		monitor.Errorf("Error: non 200 status from %s - %s", httpEndpoint, resp.Status)
		procErrChan <- httpEndpoint
		procSuccessChan <- ""
		return
	} else {
		monitor.Printf("%s returns 200 OK", httpEndpoint)
	}

	//
	// Doing this keeps the channel open
	// If this is not done, the channel closes and there is a fatal error
	//
	procErrChan <- ""
	procSuccessChan <- httpEndpoint
}

func (monitor *Monitor) checkLocalProcess(processName string, procErrChan chan string, procSuccessChan chan string, wg *sync.WaitGroup) {

	log.WithFields(log.Fields{
		"process": processName,
		"status":  STATUS_CHECK,
	}).Infof("Checking, process=%s ...", processName)

	defer wg.Done()

	pid, _, _ := findProcess(processName)

	if pid == 0 {

		log.WithFields(log.Fields{
			"process": processName,
			"status":  STATUS_CHECK_DOWN,
		}).Errorf("Error, no process=%s found running.", processName)

		procErrChan <- processName
		procSuccessChan <- ""
		return
	}

	//
	// Doing this keeps the channel open
	// If this is not done, the channel closes and there is a fatal error
	//
	procErrChan <- ""
	procSuccessChan <- processName
}
