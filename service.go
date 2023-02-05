package goservice

import (
	"errors"
	"flag"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sevlyar/go-daemon"
	"github.com/soellman/pidfile"
)

// Service main function
//
// argument string is the path to a config file supplied in startup arguments
// argument logger is a log.Logger object, if run as a service logs to syslog
type MainFunction func(string, *log.Logger)

func sMain(serviceMain MainFunction, context *daemon.Context, pidPath string, cfgPath string, useSyslog bool) {
	var logger *log.Logger = log.Default()

	if context != nil {
		defer context.Release()

		// setup pid file
		err := pidfile.Write(pidPath)
		if err != nil {
			log.Fatal(err)
		}
		defer pidfile.Remove(pidPath)
	}

	// setup logging
	// TODO: configurable log level
	slog, err := syslog.NewLogger(syslog.LOG_INFO, int(syslog.LOG_DAEMON))
	if err != nil {
		log.Fatal(err)
	}
	logger = slog // replace default logger

	// spool up service main function
	go serviceMain(cfgPath, logger)

	// wait for signal from os to quit etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	<-sigchan

}

// daemon start, supply a go-service.MainFunction() in place of a normal main() function
//
// parses startup arguments
//
//   - pidfile=/path/to/pidfile
//     if not supplied runs like a normal executable
//
//   - cfgfile=/path/to/cfgfile
//     this gets passed back to serviceMain for the user to config their service
//     function also checks that file exists if a non-empty string is passed and
//     throws a fatal error otherwise
func Start(serviceMain MainFunction) {

	// handle startup argument parsing
	var pidFilePath string
	var cfgFilePath string
	var useSyslog bool
	flag.StringVar(&pidFilePath, "pidfile", "", "supply path for pid file to run in daemon mode")
	flag.StringVar(&cfgFilePath, "cfgfile", "", "path to config file")
	flag.BoolVar(&useSyslog, "syslog", false, "set true to log to system log")
	flag.Parse()

	// if cfgFilePath is supplied, check that it exists, otherwise fatal error
	if cfgFilePath != "" {
		if _, err := os.Stat(cfgFilePath); errors.Is(err, os.ErrNotExist) {
			log.Fatal("could not find config file at", cfgFilePath)
		}
	}

	if pidFilePath != "" {
		context := new(daemon.Context)
		child, err := context.Reborn()
		if err != nil {
			log.Fatal(err)
		}
		if child == nil {
			// we're in the child process, continue
			sMain(serviceMain, context, pidFilePath, cfgFilePath, useSyslog)
		}
	} else {
		// just run the service
		sMain(serviceMain, nil, pidFilePath, cfgFilePath, useSyslog)
	}

}
