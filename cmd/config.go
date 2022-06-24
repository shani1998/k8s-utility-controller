package main

import (
	"github.com/spf13/pflag"
)

const (
	defaultServerAddr = "127.0.0.1"
	defaultServerPort = "8080"

	defaultLogLevel = "info"
)

var (
	_ = pflag.String("server.host", defaultServerAddr, "address on which server will run")
	_ = pflag.String("server.port", defaultServerPort, "port to bind the server listener to")

	_ = pflag.String("log.level", defaultLogLevel, "set the logging level(debug, info, warning, error, fatal, panic) default: info")
)
