package main

import (
	"context"

	pb "github.com/shani1998/k8s-utility-controller/proto"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	logrusEntry := log.NewEntry(log.StandardLogger())

	// Add middleware to grpc svc
	unaryInterceptorChain := []grpc.UnaryServerInterceptor{
		// logging middle interceptor for grpc logging
		grpcLogrus.UnaryServerInterceptor(logrusEntry),
		// prometheus' metrics interceptor for getting grpc metrics
		grpcPrometheus.UnaryServerInterceptor,
		// validation interceptor  for validating incoming requests
		grpcValidator.UnaryServerInterceptor(),
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(unaryInterceptorChain...)),
	)

	return grpcServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}
