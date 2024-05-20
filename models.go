package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/rxmeez/go-blog-agg/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	var lastFetchedAt *time.Time
	if feed.LastFetchedAt.Valid {
		lastFetchedAt = &feed.LastFetchedAt.Time
	}
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: lastFetchedAt,
	}
}

func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	var convertedFeeds []Feed
	for _, dbFeed := range feeds {
		var lastFetchedAt *time.Time
		if dbFeed.LastFetchedAt.Valid {
			lastFetchedAt = &dbFeed.LastFetchedAt.Time
		}
		f := Feed{
			ID:            dbFeed.ID,
			CreatedAt:     dbFeed.CreatedAt,
			UpdatedAt:     dbFeed.UpdatedAt,
			Name:          dbFeed.Name,
			Url:           dbFeed.Url,
			UserID:        dbFeed.UserID,
			LastFetchedAt: lastFetchedAt,
		}
		convertedFeeds = append(convertedFeeds, f)
	}
	return convertedFeeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
	}
}

func databaseFeedFollowsToFeedFollows(feeds []database.FeedFollow) []FeedFollow {
	var convertedFeeds []FeedFollow
	for _, dbFeed := range feeds {
		f := FeedFollow{
			ID:        dbFeed.ID,
			CreatedAt: dbFeed.CreatedAt,
			UpdatedAt: dbFeed.UpdatedAt,
			FeedID:    dbFeed.FeedID,
			UserID:    dbFeed.UserID,
		}
		convertedFeeds = append(convertedFeeds, f)
	}
	return convertedFeeds
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(post database.Post) Post {
	var description string
	if post.Description.Valid {
		description = post.Description.String
	}
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: description,
		PublishedAt: post.PublishedAt,
		FeedID:      post.FeedID,
	}
}

func databasePostsToPosts(posts []database.Post) []Post {
	var convertedPosts []Post
	for _, dbPost := range posts {
		var description string
		if dbPost.Description.Valid {
			description = dbPost.Description.String
		}
		p := Post{
			ID:          dbPost.ID,
			CreatedAt:   dbPost.CreatedAt,
			UpdatedAt:   dbPost.UpdatedAt,
			Title:       dbPost.Title,
			Url:         dbPost.Url,
			Description: description,
			PublishedAt: dbPost.PublishedAt,
			FeedID:      dbPost.FeedID,
		}
		convertedPosts = append(convertedPosts, p)
	}
	return convertedPosts
}
