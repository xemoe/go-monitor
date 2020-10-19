#!/bin/sh

set -euo pipefail

echo "Installing for linux"

cp ./bin/go-monitor /usr/local/bin/go-monitor
chmod +x /usr/local/bin/go-monitor

test ! -d /etc/g-monitor && mkdir /etc/go-monitor || true
test ! -f /etc/g-monitor/config.yml && cp ../config.yml.default /etc/go-monitor/ || true

cp gomonitor.service /etc/systemd/system/

systemctl enable gomonitor
