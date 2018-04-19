package main

import (
	"log"

	"google.golang.org/grpc"

	pubpb "github.com/pavelnikolov/eventsourcing-go/publishing"
	eventspb "github.com/pavelnikolov/eventsourcing-go/publishing/events"
	"github.com/pavelnikolov/eventsourcing-go/services/graph"
)

const (
	articlesAddress = "localhost:50051"
	gatewayAddress  = "localhost:50052"
)

func main() {
	articlesConn := connect(articlesAddress)
	defer articlesConn.Close()
	articlesClient := pubpb.NewArticlesClient(articlesConn)

	gatewayConn := connect(gatewayAddress)
	defer gatewayConn.Close()
	gatewayClient := eventspb.NewGatewayClient(gatewayConn)

	graph.StartServer(articlesClient, gatewayClient)
}

func connect(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(articlesAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	return conn
}
