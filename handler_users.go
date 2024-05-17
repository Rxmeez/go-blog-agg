package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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
		return
	}

	id := uuid.New()
	currentTime := time.Now().UTC()

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        id,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))

}

func (cfg apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
