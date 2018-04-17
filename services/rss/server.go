package rss

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/feeds"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
)

const port = "4002"

// StartServer starts http server and exposes RSS feed endpoints
func StartServer(c pb.ArticlesClient) {
	http.Handle("/feed", rssHanlder(c, ""))
	http.Handle("/feed/business", rssHanlder(c, "business"))
	http.Handle("/feed/politics", rssHanlder(c, "politics"))

	log.Printf("Listening for connections on http://localhost:%s/feed\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func rssHanlder(c pb.ArticlesClient, category string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &pb.LatestArticlesRequest{
			Category: category,
			Count:    20,
			Status:   pb.ArticleStatus_PUBLISHED,
		}
		res, err := c.LatestArticles(r.Context(), req)
		if err != nil {
			http.Error(w, "failed to query articles", http.StatusInternalServerError)
			log.Printf("failed to fetch articles: %v\n", err)
			return
		}

		feed, err := generateFeed(res.Articles)
		if err != nil {
			http.Error(w, "failed to generate RSS feed", http.StatusInternalServerError)
			log.Printf("failed to generate RSS feed: %v\n", err)
			return
		}

		w.Header().Add("Content-Type", "application/rss+xml")
		feed.WriteRss(w)
	}
}

func generateFeed(articles []*pb.Article) (*feeds.Feed, error) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "Company Name Here",
		Link:        &feeds.Link{Href: "https://example.com"},
		Description: "When news breaks, we fix it!",
		Author:      &feeds.Author{Name: "Company Name Here", Email: "contact@example.com"},
		Created:     now,
	}

	for _, a := range articles {
		created, err := ptypes.Timestamp(a.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to convert date: %v", err)
		}

		item := &feeds.Item{
			Title:       a.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("http://example.com/%s/%s-%d", toURLPath(a.Category), toURLPath(a.Title), a.Id)},
			Description: a.Body,
			Author:      &feeds.Author{Name: a.AuthorName, Email: toURLPath(a.AuthorName) + "@example.com"},
			Created:     created,
		}
		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}

func toURLPath(s string) string {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	return strings.ToLower(strings.Join(strings.FieldsFunc(s, f), "-"))
}
