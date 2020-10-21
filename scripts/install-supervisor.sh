#!/bin/sh

set -euo pipefail

#############################
# START HEADER SCRIPT #
#######################

ROOT_UID="0"
MUID=$(/usr/bin/id -u `whoami`)
CHECK_SUDO="true"

#############################
# CHECK IS ROOT UID #
#####################

if [ "true" == "${CHECK_SUDO}" ] && [ "$MUID" -ne "$ROOT_UID" ]
then
    echo " * You must be root to do that!";
    exit 1;
fi

#############################
# INSTALL REALPATH #
# IF NOT EXIST #####
####################

if ! type realpath &>/dev/null
then
    sudo apt-get update && sudo apt-get install realpath -y;
fi

SCRIPT_PATH=$(dirname $(realpath -s $0))

#############################
# BEGIN #
#########

echo "Installing with supervisor"

cp ${SCRIPT_PATH}/../bin/go-monitor /usr/local/bin/go-monitor
chmod +x /usr/local/bin/go-monitor

test ! -d /etc/g-monitor && mkdir /etc/go-monitor || true
test ! -f /etc/g-monitor/config.yml && cp ${SCRIPT_PATH}/../config.yml.default /etc/go-monitor/config.yml || true

cp ${SCRIPT_PATH}/../init/supervisor_gomonitor.conf /etc/supervisor/conf.d/supervisor_gomonitor.conf

supervisorctl reload
