package scraper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
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
						var description sql.NullString
						if item.Description != "" {
							description.String = item.Description
							description.Valid = true
						}

						pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
						if err != nil {
							pubDate, err = time.Parse(time.RFC1123, item.PubDate)
							if err != nil {
								fmt.Println("Error parsing pubDate:", err)
								return
							}
						}

						post, err := db.CreatePost(ctx, database.CreatePostParams{
							ID:          uuid.New(),
							CreatedAt:   time.Now().UTC(),
							UpdatedAt:   time.Now().UTC(),
							Title:       item.Title,
							Url:         item.Link,
							Description: description,
							PublishedAt: pubDate,
							FeedID:      feed.ID,
						})

						if err != nil {
							if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
								continue
							} else {
								fmt.Println("Error saving post to database:", err)
							}
						}

						fmt.Println("Saved post to database:", post.Url)
					}
				}(feed)
			}

			wg.Wait()
			fmt.Println("Finished fetching feeds.")

		}
	}
}
