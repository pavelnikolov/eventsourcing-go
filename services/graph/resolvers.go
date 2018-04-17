package graph

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
)

const (
	articleKind = "article"
	authorKind  = "author"
)

type queryResolver struct {
	client pb.ArticlesClient
}

func (r *queryResolver) Article(ctx context.Context, args struct{ ID graphql.ID }) (*articleResolver, error) {
	var aid int32
	relay.UnmarshalSpec(args.ID, &aid)
	a, err := r.client.Article(ctx, &pb.ArticleRequest{Id: uint32(aid)})
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

func (r *queryResolver) Articles(ctx context.Context, args struct {
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
	articles, err := r.client.LatestArticles(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %v", err)
	}

	var res []*articleResolver
	for _, a := range articles.Articles {
		res = append(res, &articleResolver{article: a})
	}

	return res, nil
}
