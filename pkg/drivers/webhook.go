package drivers

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/xemoe/go-monitor/pkg/monitor"
)

type WebHookDriver struct {
}

func (wh WebHookDriver) Init() error {
	return nil
}

func (wh WebHookDriver) Alert(m *monitor.Monitor, proc string, server monitor.ServerInfo, state string) error {

	conf := m.Drivers["webhook"].(map[string]interface{})

	//
	// Available configurations
	// ** lowercase only **
	// - apiurl
	// - authToken
	// - sender
    // - skiphttpsverify: "yes|no"
	//
	apiUrl := conf["apiurl"].(string)
	authToken := conf["authtoken"].(string)
	sender := conf["sender"].(string)
	skipHttpsVerify := conf["skiphttpsverify"].(string) == "yes"

	//
	// Send text message
	//
	urlStr := apiUrl

	v := url.Values{}
	v.Set("originator", sender)
	v.Set("process", proc)
	v.Set("state", state)
	v.Set("serverHost", server.Host)
	v.Set("serverIp", server.Ip)
	rb := *strings.NewReader(v.Encode())

	m.Debugf(v.Encode())

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipHttpsVerify},
	}
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth("AccessKey", authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//
	// Make request
	//
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//
	// resp.Body is an io.ReadCloser... NewDecoder expects an io.Reader
	//
	var result map[string]interface{}
	dec := json.NewDecoder(resp.Body)

	if err := dec.Decode(&result); err != nil {
		m.Errorf("Error decode webhook JSON response - " + err.Error())
		return nil
	}

	dec.Decode(&result)
	m.Debugf("Debug webhook JSON response - %v", result)

	return nil
}
