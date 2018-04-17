package main

import (
	"log"

	"google.golang.org/grpc"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
	"github.com/pavelnikolov/eventsourcing-go/services/graph"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewArticlesClient(conn)

	graph.StartServer(c)
}
