package drivers

import (
	"fmt"

	"github.com/gregdel/pushover"
	"github.com/xemoe/go-monitor/pkg/monitor"
)

type PushoverDriver struct {
}

func (pd PushoverDriver) Init() error {
	return nil
}

func (pd PushoverDriver) Alert(m *monitor.Monitor, proc string, server monitor.ServerInfo, state string) error {

	conf := m.Drivers["pushover"].(map[string]interface{})

	recipient := conf["recipient"].(string)
	token := conf["token"].(string)
	app := pushover.New(token)

	//
	// Create a new recipient
	//
	target := pushover.NewRecipient(recipient)

	//
	// Create the message to send
	//
	out := fmt.Sprintf("%s - %s @ %s", state, proc, server.String())

	message := pushover.NewMessage(out)

	_, err := app.SendMessage(message, target)
	if err != nil {
		return err
	}

	return nil
}
