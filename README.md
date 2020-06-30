[![Go Report Card](https://goreportcard.com/badge/github.com/rob121/go-monitor)](https://goreportcard.com/report/github.com/rob121/go-monitor)

# Overview 

Simple server monitoring in Go

`go-monitor` is simple server monitoring written in Go. I was on the lookout for a tool which allowed me 
to monitor a list of services and notify me through SMS or email. Everything I found seemed a bit too complex 
and so `go-monitor` was born.


## Notifications

Driver based notifications, supports the following

* Message bird [MessageBird](https://www.messagebird.com/)
* Pushover 


## Installation/Running 

### Manual

1. `go get github.com/rob121/go-monitor`

2. Update `go-monitor.yml.sample` to include the configuration options desired. 

3. Rename config to `config.yml`.

4 ./go-monitor 

### Linux

```
1. Clone the repo
2. cd install
3. sh install.sh
4. Update configuration in /etc/gomonitor
4. systemctl start gomonitor
```

## Compatibility 

I've not tested on windows, however the process checking library now supports it, if you do test on windows, please let me know.

## Configuration

Config is read through a yaml file.
Users can specify how often they are notified of a given service being down through the `defaultttl` value in the config.

By default, processes are checked every 60 seconds. This can be increased or decreased depending on the importance of services.

```
processes: [ "http://192.168.1.210", "go-monitor"] //list of processes to check
notificationdriver: "pushover" //driver name, match the key below
config:
   defaultttlseconds: 3600  //how long to hold the "down" state of a service, shorter means more notifications
   notifyservicereturn: true //send a notification when the service returns up?
   checkfrequencyseconds: 5 //how often to check
drivers:
   messagebird:
     token: "test_TOKEN"
     sender: "+sender-number"
     recipients: "+recipient-numbers,+one-or-many"
   pushover: 
     token: "test"
     recipient: "test"
```

## Additions
Down the line it would be nice to have server monitoring in addition to process monitoring:

- CPU
- Disks (usage, read/write)
- Memory

This could further be extended into network monitoring. Someday.

## License
MIT
