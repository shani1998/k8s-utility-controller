package main

import (
	"github.com/spf13/pflag"
)

const (
	defaultServerAddr = "127.0.0.1"
	defaultServerPort = "8080"

	defaultHealthServerEnable = true
	defaultHealthAddress      = "0.0.0.0"
	defaultHealthPort         = "8089"

	defaultLogLevel  = "info"
	defaultLogFormat = "json"
)

var (
	_ = pflag.String("server.host", defaultServerAddr, "address on which server will run")
	_ = pflag.String("server.port", defaultServerPort, "port to bind the server listener to")

	_ = pflag.Bool("healthz.enable", defaultHealthServerEnable, "the flag that indicates whether the heath-check endpoint is enabled, default: true")
	_ = pflag.String("healthz.host", defaultHealthAddress, "address and port to bind the health check listener to")
	_ = pflag.String("healthz.port", defaultHealthPort, "port to bind the health check listener to")

	_ = pflag.String("log.level", defaultLogLevel, "set the logging level(debug, info, warning, error, fatal, panic) default: info")
	_ = pflag.String("log.format", defaultLogFormat, "set the logging format(json,text) default: json")
)
