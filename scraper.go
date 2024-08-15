package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Anas-Sayed0/rss-agg/internal/database"
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
        log.Println("found post", item.Title, "on feed:", feed.Name)
    }
    log.Printf("feed %s collected %v posts found", feed.Name, len(rssFeed.Channel.Items))
}

