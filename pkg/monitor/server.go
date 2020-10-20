package monitor

import (
	"net"
	"os"
)

type ServerInfo struct {
	Host string
	Ip   string
}

func (s ServerInfo) String() string {
	return s.Host + " with IP " + s.Ip
}

func (monitor *Monitor) getServerInfo() (server ServerInfo, err error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return ServerInfo{}, err
	}

	var ip net.IP
	for _, i := range ifaces {

		addrs, err := i.Addrs()
		if err != nil {
			return ServerInfo{}, err
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
		}
	}

	//
	// Get hostname
	//
	host, err := os.Hostname()
	if err != nil {
		return ServerInfo{}, err
	}

	return ServerInfo{Host: host, Ip: ip.String()}, nil

}
