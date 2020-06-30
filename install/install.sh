#!/bin/sh


echo "Installing for linux"

wget -O /usr/local/bin/  "http://github"
chmod +x /usr/local/bin/go/monitor
mkdir /etc/gomonitor
cp ../sample.config.yml /etc/gomonitor/
cp gomonitor.service /etc/systemd/system/
systemctl enable gomonitor