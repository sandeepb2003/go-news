package domain

import "time"

type Article struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Published   *time.Time `json:"published_at"`
	Link        string     `json:"link"`
	FeedTitle   string     `json:"feed_title"`
}
