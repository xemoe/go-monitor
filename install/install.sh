#!/bin/sh


echo "Installing for linux"

wget -O /usr/local/bin/go-monitor  "https://github.com/rob121/go-monitor/releases/download/v.0.1/go-monitor-linux-amd64"
chmod +x /usr/local/bin/go-monitor
mkdir /etc/gomonitor
cp ../sample.config.yml /etc/gomonitor/
cp gomonitor.service /etc/systemd/system/
systemctl enable gomonitor
