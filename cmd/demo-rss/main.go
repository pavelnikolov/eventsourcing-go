package main

import (
	"log"

	"google.golang.org/grpc"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
	"github.com/pavelnikolov/eventsourcing-go/services/rss"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewArticlesClient(conn)

	rss.StartServer(c)
}
