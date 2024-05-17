package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rxmeez/go-blog-agg/internal/database"
)

func (cfg apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}

	id := uuid.New()
	currentTime := time.Now().UTC()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate api key")
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        id,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))

}

func (cfg apiConfig) handlerAllFeedsGet(w http.ResponseWriter, r *http.Request) {
	// create sql to retrieve all feeds and call and respond with a list of feeds
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve the feed")
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
