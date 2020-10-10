package drivers

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/xemoe/go-monitor/pkg/monitor"
)

type WebHookDriver struct {
}

func (wh WebHookDriver) Init() error {
	log.Println("Init WebHook!")
	return nil
}

func (wh WebHookDriver) Alert(m *monitor.Monitor, proc string, server monitor.ServerInfo, state string) error {

	conf := m.Drivers["webhook"].(map[string]interface{})

	webhookAPI := conf["apiUrl"].(string)
	recipientNumber := conf["recipients"].(string)
	authToken := conf["token"].(string)
	sender := conf["sender"].(string)

	//
	// Send text message
	//
	urlStr := webhookAPI

	v := url.Values{}
	v.Set("recipients", recipientNumber)
	v.Set("originator", sender)
	v.Set("body", "📢 "+proc+" "+state+" on "+server.String()+"!")
	rb := *strings.NewReader(v.Encode())

	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth("AccessKey", authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//
	// Make request
	//
	_, err := client.Do(req)
	if err != nil {
		return err
	}

	m.Println("Notification sent!")

	return nil
}
