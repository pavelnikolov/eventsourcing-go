package sitemap

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/golang/protobuf/ptypes"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
)

const port = "4003"

// StartServer starts http server and exposes RSS feed endpoints
func StartServer(c pb.ArticlesClient) {
	http.Handle("/sitemap", sitemapHanlder(c))

	log.Printf("Listening for connections on http://localhost:%s/sitemap\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func sitemapHanlder(c pb.ArticlesClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &pb.LatestArticlesRequest{
			Count:  20,
			Status: pb.ArticleStatus_PUBLISHED,
		}
		res, err := c.LatestArticles(r.Context(), req)
		if err != nil {
			http.Error(w, "failed to query articles", http.StatusInternalServerError)
			log.Printf("failed to fetch articles: %v\n", err)
			return
		}

		sm := buildSitemap(res.Articles)

		w.Header().Add("Content-Type", "application/xml")
		w.Write(sm.XMLContent())
	}
}

func buildSitemap(articles []*pb.Article) *stm.Sitemap {
	sm := stm.NewSitemap()
	sm.SetDefaultHost("http://example.com")

	sm.Create()
	sm.Add(stm.URL{"loc": "/", "changefreq": "daily"})
	sm.Add(stm.URL{"loc": "/about", "mobile": true})

	var urls []stm.URL
	// generate a very naive and useless sitemap
	for _, a := range articles {
		urls = append(urls, stm.URL{
			"loc":              fmt.Sprintf("http://example.com/%s/%s-%d", toURLPath(a.Category), toURLPath(a.Title), a.Id),
			"title":            a.Title,
			"keywords":         strings.Split(a.Title, " "),
			"publication_date": ptypes.TimestampString(a.Created),
			"access":           "Subscription",
			"genres":           a.Category,
		})
	}

	sm.Add(stm.URL{"loc": "/news", "changefreq": "hourly", "news": urls})

	// Note: Do not call `sm.Finalize()` because it flushes
	// the underlying datastructure from memory to disk.
	return sm
}

func toURLPath(s string) string {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	return strings.ToLower(strings.Join(strings.FieldsFunc(s, f), "-"))
}
