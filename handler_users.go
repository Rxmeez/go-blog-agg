package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rxmeez/go-blog-agg/internal/auth"
	"github.com/rxmeez/go-blog-agg/internal/database"
)

func (cfg apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
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

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        id,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))

}

func (cfg apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect ApiKey")
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ApiKey doesn't match any records")
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
