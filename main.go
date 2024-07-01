package main

import (
	"log"
	"net"

	"github.com/IgorEulalio/notificationservice/cmd/server"
	pb "github.com/IgorEulalio/notificationservice/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

const createRepositoryTimeoutSeconds = "5"

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to establish NATS connection: %v", err)
	}

	pb.RegisterRepositoryServiceServer(s, &server.Server{Nc: nc})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
