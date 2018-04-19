//go:generate protoc -I ../../publishing/broker --go_out=plugins=grpc:../../publishing/broker ../../publishing/broker/broker.proto

package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/jzelinskie/grpc/simple"
)

type simpleServer struct {
}

func (s *simpleServer) SimpleRPC(stream pb.SimpleService_SimpleRPCServer) error {
	log.Println("Started stream")
	for {
		in, err := stream.Recv()
		log.Println("Received value")
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Got " + in.Msg)
	}
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleServiceServer(grpcServer, &simpleServer{})

	l, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Listening on tcp://localhost:6000")
	grpcServer.Serve(l)
}
