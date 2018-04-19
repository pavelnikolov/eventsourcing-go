package gateway

import (
	"errors"
	"log"

	"golang.org/x/net/context"

	brokerpb "github.com/pavelnikolov/eventsourcing-go/publishing/broker"
	eventspb "github.com/pavelnikolov/eventsourcing-go/publishing/events"
)

// package errors
var (
	ErrMissingBody     = errors.New("article body is required")
	ErrMissingCategory = errors.New("article category is required")
	ErrMissingTitle    = errors.New("article title is required")
	ErrNilArticle      = errors.New("article is <nil>")
	ErrUnknownStatus   = errors.New("unknown article status")
)

// NewServer returns a new isntance of the server.
func NewServer() *Server {
	return &Server{}
}

// Server accepts and validates events.
type Server struct {
	brokerClient *brokerpb.BrokerClient
}

// ArticleCreated receives an ArticleCreated event, validates it and forwards it to the broker
func (s *Server) ArticleCreated(ctx context.Context, in *eventspb.ArticleCreatedRequest) (*eventspb.ArticleCreatedReply, error) {
	log.Printf("Received event: %v", in.Article)
	return &eventspb.ArticleCreatedReply{}, nil
}

// ArticleUpdated receives an ArticleUpdated event, validates it and forwards it to the broker
func (s *Server) ArticleUpdated(ctx context.Context, in *eventspb.ArticleUpdatedRequest) (*eventspb.ArticleUpdatedReply, error) {
	log.Printf("Received event: %v", in.Article)
	return &eventspb.ArticleUpdatedReply{}, nil
}

// ArticleRetracted receives an ArticleRetracted event, validates it and forwards it to the broker
func (s *Server) ArticleRetracted(ctx context.Context, in *eventspb.ArticleRetractedRequest) (*eventspb.ArticleRetractedReply, error) {
	log.Printf("Received event: %v", in.Article)
	return &eventspb.ArticleRetractedReply{}, nil
}

// ArticleDraftCreated receives an ArticleDraftCreated event, validates it and forwards it to the broker
func (s *Server) ArticleDraftCreated(ctx context.Context, in *eventspb.ArticleDraftCreatedRequest) (*eventspb.ArticleDraftCreatedReply, error) {
	log.Printf("Received event: %v", in.Article)
	return &eventspb.ArticleDraftCreatedReply{}, nil
}

// ArticleDraftUpdated receives an ArticleDraftUpdated event, validates it and forwards it to the broker
func (s *Server) ArticleDraftUpdated(ctx context.Context, in *eventspb.ArticleDraftUpdatedRequest) (*eventspb.ArticleDraftUpdatedReply, error) {
	log.Printf("Received event: %v\n", in.Article)
	return &eventspb.ArticleDraftUpdatedReply{}, nil
}

func validate(a *eventspb.Article) error {
	if a == nil {
		return ErrNilArticle
	}
	if a.Body == "" {
		return ErrMissingBody
	}
	if a.Category == "" {
		return ErrMissingCategory
	}
	if a.Title == "" {
		return ErrMissingTitle
	}
	return nil
}
