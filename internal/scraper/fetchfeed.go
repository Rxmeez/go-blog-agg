package scraper

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rxmeez/go-blog-agg/internal/database"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(url string) (Rss, error) {

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return Rss{}, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return Rss{}, err
	}

	rss := Rss{}
	err = xml.Unmarshal(bodyBytes, &rss)
	if err != nil {
		fmt.Println(err)
		return Rss{}, err
	}

	return rss, nil
}

func StartScraping(ctx context.Context, db *database.Queries, interval time.Duration, n int32) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker stopped.")
			return
		case <-ticker.C:
			fmt.Println("Fetching feeds...")

			feeds, err := db.GetNextFeedsToFetch(ctx, n)
			if err != nil {
				log.Fatalln("Error fetching feed:", err)
				continue
			}

			var wg sync.WaitGroup
			for _, feed := range feeds {
				wg.Add(1)
				go func(feed database.GetNextFeedsToFetchRow) {
					defer wg.Done()

					rss, err := FetchFeed(feed.Url)
					if err != nil {
						fmt.Println("Error fetching feed:", err)
						return
					}

					for _, item := range rss.Channel.Items {
						fmt.Println("Title:", item.Title)
					}
				}(feed)
			}

			wg.Wait()
			fmt.Println("Finished fetching feeds.")

		}
	}
}
