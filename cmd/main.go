package main

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/shani1998/k8s-utility-controller/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	initializeLogger(viper.GetString("log.format"), viper.GetString("log.level"))
}

func main() {
	// setup check endpoint to monitor the health of the controller,
	// it becomes unhealthy if at all any error occurs while serving requests
	if viper.GetBool("healthz.enable") {
		go updateHealthStatus()
		go healthz(viper.GetString("healthz.host"), viper.GetString("healthz.port"))
	}

	// retry until we initialize the kube client successfully
	// using either inCluster or outCluster config.
	for handlers.InitKubeClient() != nil {
		// Wait for 30s to 1m before making a request to api server
		jitter := time.Duration(rand.Intn(30*1000)) * time.Millisecond
		duration := 30*time.Second + jitter
		log.Infof("retry initializing kube client in %v", duration)
		time.Sleep(duration)
	}

	// initialize http router
	router := httprouter.New()
	// get services
	router.GET("/services", handlers.GetServices)
	// get services by application group
	router.GET("/services/:applicationGroup", handlers.GetServicesByAppLabel)

	srv := &http.Server{
		Addr:    net.JoinHostPort(viper.GetString("server.host"), viper.GetString("server.port")),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error while servingon %v", err)
		}
	}()
	log.Infof("started server on address %s", srv.Addr)

	// accepts os Interrupt signal to shut down server gracefully
	doneC := make(chan os.Signal, 1)
	signal.Notify(doneC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Infof("received signal %v", <-doneC)

	// shut down http server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful server shutdown failed: %v", err)
	}
	log.Infof("gracefully stopped server listening on %s", srv.Addr)
}
