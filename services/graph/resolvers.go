package graph

import (
	"context"
	"errors"
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
	eventspb "github.com/pavelnikolov/eventsourcing-go/publishing/events"
)

const (
	articleKind = "article"
	authorKind  = "author"
)

type resolver struct {
	articlesClient pb.ArticlesClient
	gatewayClient  eventspb.GatewayClient
}

func (r *resolver) Article(ctx context.Context, args struct{ ID graphql.ID }) (*articleResolver, error) {
	var aid int32
	relay.UnmarshalSpec(args.ID, &aid)
	a, err := r.articlesClient.Article(ctx, &pb.ArticleRequest{Id: uint32(aid)})
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %v", err)
	}

	return &articleResolver{article: a.Article}, nil
}

type articleResolver struct {
	article *pb.Article
}

func (r *articleResolver) ID() graphql.ID {
	return relay.MarshalID(articleKind, r.article.Id)
}

func (r *articleResolver) Body() string {
	return r.article.Body
}

func (r *articleResolver) Category() string {
	return r.article.Category
}

func (r *articleResolver) Title() string {
	return r.article.Title
}

func (r *articleResolver) AuthorID() graphql.ID {
	return relay.MarshalID(authorKind, r.article.AuthorId)
}

func (r *articleResolver) AuthorName() string {
	return r.article.AuthorName
}

func (r *articleResolver) Status() string {
	return r.article.Status.String()
}

func (r *resolver) Articles(ctx context.Context, args struct {
	Category *string
	Count    int32
	Status   string
}) ([]*articleResolver, error) {
	var category string
	if args.Category != nil {
		category = *args.Category
	}
	req := &pb.LatestArticlesRequest{
		Category: category,
		Count:    uint32(args.Count),
		Status:   pb.ArticleStatus(pb.ArticleStatus_value[args.Status]),
	}
	articles, err := r.articlesClient.LatestArticles(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %v", err)
	}

	var res []*articleResolver
	for _, a := range articles.Articles {
		res = append(res, &articleResolver{article: a})
	}

	return res, nil
}

type articleCreatedInput struct {
	Article *articleInput
}

func (r *resolver) CreateArticle(ctx context.Context, args articleCreatedInput) (*eventResponseResolver, error) {
	status := "FAILED"
	err := validateArticleInput(args.Article)
	if err != nil {
		return &eventResponseResolver{
			id:     relay.MarshalID(articleKind, args.Article.ID),
			status: status,
			err:    err,
		}, nil
	}

	req := &eventspb.ArticleCreatedRequest{
		Article: &eventspb.Article{
			Id:         uint32(args.Article.ID),
			Title:      args.Article.Title,
			Body:       args.Article.Body,
			Category:   args.Article.Category,
			AuthorId:   uint32(args.Article.AuthorID),
			AuthorName: args.Article.AuthorName,
		},
	}

	_, err = r.gatewayClient.ArticleCreated(ctx, req)
	if err != nil {
		return &eventResponseResolver{
			id:     relay.MarshalID(articleKind, args.Article.ID),
			status: status,
			err:    fmt.Errorf("failed to send event to the gateway: %v", err),
		}, nil
	}

	status = "ACCEPTED"
	return &eventResponseResolver{
		id:     relay.MarshalID(articleKind, args.Article.ID),
		status: status,
	}, nil
}

func validateArticleInput(a *articleInput) error {
	if a.Title == "" {
		return errors.New("Missing title")
	}
	if a.Category == "" {
		return errors.New("Missing category")
	}
	if a.AuthorName == "" {
		return errors.New("Missing author name")
	}
	if a.AuthorID == 0 {
		return errors.New("Missing author ID")
	}
	if a.ID == 0 {
		return errors.New("Missing article ID")
	}
	return nil
}

type articleInput struct {
	ID         int32
	Title      string
	Body       string
	Category   string
	AuthorID   int32
	AuthorName string
}

type eventResponseResolver struct {
	id     graphql.ID
	name   string
	err    error
	status string
}

func (r *eventResponseResolver) ID() graphql.ID {
	return r.id
}

func (r *eventResponseResolver) Name() string {
	return r.name
}

func (r *eventResponseResolver) Error() *string {
	if r.err != nil {
		err := r.err.Error()
		return &err
	}
	return nil
}
func (r *eventResponseResolver) Status() string {
	return r.status
}
