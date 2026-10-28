package main

import (
    "fmt"
    "os"

    "github.com/sandeepb2003/go-news/internal/config"
    "github.com/sandeepb2003/go-news/internal/news"
)

func main() {
    cfg := config.Default()

    if port := os.Getenv("PORT"); port != "" {
        cfg.Port = port
    }

    service := news.NewService(cfg.FeedURL)

    fmt.Printf("go-news initialized successfully\n")
    fmt.Printf("Port: %s\n", cfg.Port)
    fmt.Printf("Feed URL: %s\n", cfg.FeedURL)
    fmt.Printf("Sample headlines: %v\n", service.Headlines())
}
