# go-service
Package to quickly implement linux services in go

## why?
I've written this in one form or another a good number of times and I'd rather not continue doing so.

## how to use
TODO

```go
func serviceMain(cfgFilePath string, serviceLogger *log.Logger) {
    // do stuff
}

func main() {
    // let goservice handle the leg work
	goservice.Start(serviceMain)
}

```
### works with systemd
example:
```
[Unit]
Description=Go Service Example
After=network.target

[Service]
Type=simple
User=some_service_user
Group=some_service_group
PIDFile=/var/run/my_exe.pid
ExecStart=/path/to/my_exe --cfgfile=/path/to/my/config_file --pidfile=/var/run/my_exe.pid

# Give a reasonable amount of time for the server to start up/shut down
TimeoutSec=300

Restart=on-failure
RestartPreventExitStatus=1

# Sets open_files_limit
LimitNOFILE = 10000

[Install]
WantedBy=multi-user.target
```

### works with openrc too
```
#!/sbin/openrc-run
  
name=$RC_SVCNAME
cfgfile="/path/to/my.conf"
pidfile="/run/$RC_SVCNAME/$RC_SVCNAME.pid"

command="/usr/bin/my_service"
command_args="--pidfile=$pidfile --cfgfile=$cfgfile"
command_user="service_user"
command_background="yes"

depend() {
        need net
}

```