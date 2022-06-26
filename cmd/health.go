package main

import (
	"net"
	"net/http"

	"github.com/shani1998/k8s-utility-controller/handlers"
	log "github.com/sirupsen/logrus"
)

var healthError error

// healthz starts http server which handles the requests for readiness check
func healthz(host, port string) {
	log.Infof("starting healthz at %s:%s", host, port)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if healthError != nil {
			log.Errorf("health check failed, error: %v", healthError)
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}
	})
	address := net.JoinHostPort(host, port)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Errorf("failed to start healthz server, Reason %v", err)
	}
}

// updateHealthStatus updates the values of healthError variable upon
// receiving the error value  on channel from handlers.
func updateHealthStatus() {
	for {
		healthError = <-handlers.HealthChan
		if healthError != nil {
			log.Debugf("received health error msg from channel, %v", healthError)
		}
	}
}
