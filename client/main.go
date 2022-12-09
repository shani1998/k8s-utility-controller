// Package main implements a client for hello world service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/shani1998/k8s-utility-controller/proto"

	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := connect(context.Background(), *addr)
	if err != nil {
		log.Fatalf("unable to connect to server: %v", err)
	}
	defer conn.Close()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cl := pb.NewGreeterClient(conn)
	r, err := cl.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}

func connect(ctx context.Context, uri string) (*grpc.ClientConn, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // client side load balancing
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(grpcprometheus.UnaryClientInterceptor),
	}
	conn, err := grpc.DialContext(ctx, uri, dialOpts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
