package articles

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
)

// package errors
var (
	ErrMissingBody     = errors.New("article body is required")
	ErrMissingCategory = errors.New("article category is required")
	ErrMissingTitle    = errors.New("article title is required")
	ErrNilArticle      = errors.New("article is <nil>")
	ErrUnknownStatus   = errors.New("unknown article status")
)

// Factory is the interface of data store for articles.
type Factory interface {
	Get(ctx context.Context, id uint32) (*pb.Article, error)
	Create(ctx context.Context, a *pb.Article) (*pb.Article, error)
	Update(ctx context.Context, a *pb.Article) (*pb.Article, error)
	Latest(ctx context.Context, category string, count uint32, status pb.ArticleStatus) ([]*pb.Article, error)
}

// NewServer initialises an instance of the articles server.
func NewServer(db Factory) *Server {
	if db == nil {
		panic("db cannot be <nil>.")
	}
	return &Server{db: db}
}

// Server is used to implement publising.ArticlesServer.
type Server struct {
	db Factory
}

// Article returns an article by ID.
func (s *Server) Article(ctx context.Context, in *pb.ArticleRequest) (*pb.ArticleReply, error) {
	a, err := s.db.Get(ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get article: %v", err))
	}
	return &pb.ArticleReply{Article: a}, nil
}

// CreateArticle creates an article.
func (s *Server) CreateArticle(ctx context.Context, in *pb.CreateArticleRequest) (*pb.ArticleReply, error) {
	if err := validate(in.Article); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invslid input: %v", err))
	}

	a, err := s.db.Create(ctx, in.Article)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create article: %v", err))
	}

	return &pb.ArticleReply{Article: a}, nil
}

// UpdateArticle updates existing article.
func (s *Server) UpdateArticle(ctx context.Context, in *pb.UpdateArticleRequest) (*pb.ArticleReply, error) {
	if err := validate(in.Article); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid input: %v", err))
	}

	a, err := s.db.Update(ctx, in.Article)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update article: %v", err))
	}

	return &pb.ArticleReply{Article: a}, nil
}

// LatestArticles queries for latest articles by the given params.
func (s *Server) LatestArticles(ctx context.Context, in *pb.LatestArticlesRequest) (*pb.ArticlesReply, error) {
	if in.Count == 0 {
		return nil, status.Error(codes.InvalidArgument, "count cannot be 0")
	}
	if in.Count > 50 {
		return nil, status.Error(codes.InvalidArgument, "count cannot be greater than 50")
	}
	res, err := s.db.Latest(ctx, in.Category, in.Count, in.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get latest articles: %v", err))
	}

	return &pb.ArticlesReply{Articles: res}, nil
}

func validate(a *pb.Article) error {
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
	if a.Status == pb.ArticleStatus_UNKNOWN {
		return ErrUnknownStatus
	}
	return nil
}
