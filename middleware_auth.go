package main

import (
	"net/http"

	"github.com/rxmeez/go-blog-agg/internal/auth"
	"github.com/rxmeez/go-blog-agg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect ApiKey")
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ApiKey doesn't match any records")
		}

		handler(w, r, user)
	}
}
