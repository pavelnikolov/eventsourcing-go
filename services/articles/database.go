package articles

import (
	"errors"
	"sync"

	pb "github.com/pavelnikolov/eventsourcing-go/publishing"
	"golang.org/x/net/context"
)

// errors
var (
	ErrArticleNotFound = errors.New("article not found")
)

// Database simulates database wrapper component
type Database struct {
	data []*pb.Article
	sync.RWMutex
}

// Get returns an article by ID
func (d *Database) Get(ctx context.Context, id uint32) (*pb.Article, error) {
	d.RLock()
	defer d.RUnlock()

	for _, a := range d.data {
		if a.Id == id {
			return a, nil
		}
	}
	return nil, ErrArticleNotFound
}

// Create creates an article
func (d *Database) Create(ctx context.Context, a *pb.Article) (*pb.Article, error) {
	d.Lock()
	defer d.Unlock()

	d.data = append([]*pb.Article{a}, d.data...)
	return a, nil
}

// Update checks if an article exists in the data store and modifies it
func (d *Database) Update(ctx context.Context, a *pb.Article) (*pb.Article, error) {
	d.Lock()
	defer d.Unlock()

	for i := range d.data {
		if a.Id == d.data[i].Id {
			d.data[i] = a
			return a, nil
		}
	}

	return nil, ErrArticleNotFound
}

// Latest returns the latest articles from the data store filtered by category and status.
func (d *Database) Latest(ctx context.Context, category string, count uint32, status pb.ArticleStatus) ([]*pb.Article, error) {
	d.RLock()
	defer d.RUnlock()

	var res []*pb.Article
	var i uint32
	for _, a := range d.data {
		if (category == "" || a.Category == category) && (status == pb.ArticleStatus_UNKNOWN || a.Status == status) {
			res = append(res, a)
			i++
		}
		if i == count {
			break
		}
	}

	return res, nil
}
