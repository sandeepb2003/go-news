package domain

type Feed struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Link        string     `json:"url"`
	Articles    []*Article `json:"articles"`
}
