#!/bin/sh

echo "Installing for linux"

wget -O /usr/local/bin/go-monitor  "https://github.com/xemoe/go-monitor/releases/download/v.0.1/go-monitor-linux-amd64"
chmod +x /usr/local/bin/go-monitor
mkdir /etc/go-monitor

cp ../config.yml.default /etc/go-monitor/
cp gomonitor.service /etc/systemd/system/

systemctl enable gomonitor
