package main

import (
	"net/http"

	"github.com/rxmeez/go-blog-agg/internal/database"
)

func (cfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := cfg.DB.GetPostByUser(r.Context(), database.GetPostByUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve feed follow data")
	}

	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
