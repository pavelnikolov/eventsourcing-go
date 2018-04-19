//go:generate protoc -I ../../publishing/events --go_out=plugins=grpc:../../publishing/events ../../publishing/events/events.proto

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing/events"
	"github.com/pavelnikolov/eventsourcing-go/services/gateway"
)

const (
	port = ":50052"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterGatewayServer(s, gateway.NewServer())
	reflection.Register(s)

	log.Printf("listening for gRPC connections on localhost%v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
