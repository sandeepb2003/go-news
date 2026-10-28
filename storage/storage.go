package storage

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sandeepb2003/go-news/domain"
	"go.etcd.io/bbolt"
)

const bucketName = "articles"

var _ domain.Storage = (*BoltStore)(nil)

type BoltStore struct {
	db *bolt.DB
}

func NewBoltStore(path string, readOnly bool) (*BoltStore, error) {
	db, err := bolt.Open(path, 0600, &bbolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: readOnly,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %w", err)
	}

	if !readOnly {
		if err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			return err
		}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &BoltStore{db: db}, nil
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}

func (s *BoltStore) AddArticles(articles []*domain.Article) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("articles bucket not found")
		}
		for _, article := range articles {
			data, err := json.Marshal(article)
			if err != nil {
				return fmt.Errorf("failed to marshal article: %w", err)
			}
			if err := bucket.Put([]byte(article.ID), data); err != nil {
				return fmt.Errorf("failed to put article: %w", err)
			}
		}
		return nil
	})
}

func (s *BoltStore) GetRecent(n int) []*domain.Article {
	articles := make([]*domain.Article, 0, n)
	if n <= 0 {
		return articles
	}

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("articles bucket not found")
		}

		bucket.ForEach(func(k, v []byte) error {
			var article domain.Article
			if err := json.Unmarshal(v, &article); err != nil {
				return fmt.Errorf("failed to unmarshal article: %w", err)
			}
			articles = append(articles, &article)
			return nil
		})
		return nil
	})
	if err != nil {
		fmt.Println("Error getting recent articles:", err)
	}
	slices.SortFunc(articles, func(a, b *domain.Article) int {
		if a.Published == nil && b.Published == nil {
			return 0
		}
		if a.Published == nil {
			return 1
		}
		if b.Published == nil {
			return -1
		}
		return b.Published.Compare(*a.Published)
	})

	if n > len(articles) {
		n = len(articles)
	}
	return articles[:n]
}

func (s *BoltStore) GetByID(id string) (*domain.Article, error) {
	var article domain.Article
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("articles bucket not found")
		}
		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("article not found")
		}
		if err := json.Unmarshal(data, &article); err != nil {
			return fmt.Errorf("failed to unmarshal article: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &article, nil
}
