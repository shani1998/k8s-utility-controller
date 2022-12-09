// Package main implements a server for hello world service.
package main

import (
	"flag"
	"fmt"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	pb "github.com/shani1998/k8s-utility-controller/proto"
	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func main() {
	flag.Parse()

	//doneC := make(chan error)
	//globalCtx, globalCancel := context.WithCancel(context.Background())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := NewGRPCServer()
	pb.RegisterGreeterServer(s, &server{})

	go func() {
		metricsMux := http.NewServeMux()
		// After all your registrations, make sure all of the Prometheus metrics are initialized.
		grpc_prometheus.Register(s)
		// Register Prometheus metrics handler.
		metricsMux.Handle("/metrics", promhttp.Handler())
		address := net.JoinHostPort("", "9090")
		log.Infof("Starting metric server at address [%s]", address)
		if err := http.ListenAndServe(address, metricsMux); err != nil {
			log.Errorf("Failed to serve metrics at [%s]: %v", address, err)
		}
	}()

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
