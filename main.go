package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/shani1998/k8s-utility-controller/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {
	// retry until we initialize the kube client successfully
	// using service token of pod.
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

	doneC := make(chan os.Signal, 1)
	signal.Notify(doneC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("Server Started")

	<-doneC
	log.Infof("Server Stopped")

	// shut down http server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Infof("Server Exited Properly")
}
