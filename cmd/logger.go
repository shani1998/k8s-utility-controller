package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	filePathPrefix = "/go/src/github.com/shani1998/k8s-utility-controller/"
)

var formatterMap = map[string]logrus.Formatter{
	"json": &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(fr *runtime.Frame) (function string, file string) {
			file = fmt.Sprintf("%s:%d", strings.TrimPrefix(fr.File, filePathPrefix), fr.Line)
			return
		},
	},
	"text": &logrus.TextFormatter{},
}

func initializeLogger(format, level string) {
	formatter, ok := formatterMap[format]
	if !ok {
		logrus.Errorf("invalid log format [%v]", format)
	}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Errorf("invalid log level [%v], error: %v", logLevel, err)
	}
	logrus.SetFormatter(formatter)
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)
}
