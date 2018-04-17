//go:generate protoc -I ../../publishing --go_out=plugins=grpc:../../publishing ../../publishing/publishing.proto

package main

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
	"github.com/pavelnikolov/eventsourcing-go/services/articles"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	db := &articles.Database{}
	populateContent(db)

	pb.RegisterArticlesServer(s, articles.NewServer(db))
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func populateContent(db *articles.Database) {

	articles := []*pb.Article{
		{
			Id:         1,
			Title:      "My article title 1",
			Body:       "some articl text here 1",
			Category:   "business",
			AuthorId:   10,
			AuthorName: "Pavel",
			Status:     pb.ArticleStatus_PUBLISHED,
			Created:    ptypes.TimestampNow(),
		},
		{
			Id:         2,
			Title:      "My article title 2",
			Body:       "some articl text here 2",
			Category:   "politics",
			AuthorId:   10,
			AuthorName: "Someone",
			Status:     pb.ArticleStatus_PUBLISHED,
			Created:    ptypes.TimestampNow(),
		},
		{
			Id:         3,
			Title:      "My article title 3",
			Body:       "some articl text here 3",
			Category:   "business",
			AuthorId:   10,
			AuthorName: "Alicia G.",
			Status:     pb.ArticleStatus_PUBLISHED,
			Created:    ptypes.TimestampNow(),
		},
		{
			Id:         4,
			Title:      "My article title 4",
			Body:       "some articl text here 4",
			Category:   "lifestyle",
			AuthorId:   10,
			AuthorName: "Peter Pan",
			Status:     pb.ArticleStatus_DRAFT,
			Created:    ptypes.TimestampNow(),
		},
		{
			Id:         5,
			Title:      "My article title 5",
			Body:       "some articl text here 5",
			Category:   "lifestyle",
			AuthorId:   10,
			AuthorName: "Peter H.",
			Status:     pb.ArticleStatus_PUBLISHED,
			Created:    ptypes.TimestampNow(),
		},
		{
			Id:         6,
			Title:      "My article title 6",
			Body:       "some articl text here 6",
			Category:   "environment",
			AuthorId:   10,
			AuthorName: "John Smith",
			Status:     pb.ArticleStatus_RETRACTED,
			Created:    ptypes.TimestampNow(),
		},
	}

	for _, a := range articles {
		db.Create(context.Background(), a)
	}
}
