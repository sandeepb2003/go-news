package config

type Config struct {
    Port    string
    FeedURL string
}

func Default() Config {
    return Config{
        Port:    "8080",
        FeedURL: "https://example.com/rss",
    }
}
