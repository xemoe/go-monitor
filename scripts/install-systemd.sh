#!/bin/sh

set -euo pipefail

if ! type realpath &>/dev/null
then
    sudo apt-get update && sudo apt-get install realpath -y;
fi

SCRIPT_PATH=$(dirname $(realpath -s $0))

echo "Installing with systemd"

cp ${SCRIPT_PATH}/../bin/go-monitor /usr/local/bin/go-monitor
chmod +x /usr/local/bin/go-monitor

test ! -d /etc/g-monitor && mkdir /etc/go-monitor || true
test ! -f /etc/g-monitor/config.yml && cp ${SCRIPT_PATH}/../config.yml.default /etc/go-monitor/ || true

cp ${SCRIPT_PATH}/../init/gomonitor.service /etc/systemd/system/

systemctl enable gomonitor
