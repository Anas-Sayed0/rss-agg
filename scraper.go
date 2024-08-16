package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Anas-Sayed0/rss-agg/internal/database"
	"github.com/google/uuid"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
    log.Printf("Scraping on %v goroutines every %v durations", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err:= db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
    defer wg.Done()

    _, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
    if err != nil {
        log.Println("error marking feed as fetched:", err)
        return
    }

    rssFeed, err := urlToFeed(feed.Url)
    if err != nil {
        log.Printf("error fetching feed %s: %v", feed.Url, err)
        return
    }

    if len(rssFeed.Channel.Items) == 0 {
        log.Printf("feed %s collected 0 posts found", feed.Name)
        return
    }

    for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
	
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("error parsing time %s: %v", item.PubDate, err)
			continue
		}
	
		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			})
		if err != nil {
			if strings.Contains(err.Error(),"duplicate key") {
				continue
			}
			log.Printf("error creating post %s: %v", item.Title, err)
		}
	}
	log.Printf("feed %s collected %v posts found", feed.Name, len(rssFeed.Channel.Items))
}

