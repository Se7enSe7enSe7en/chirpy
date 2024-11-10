package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
	"github.com/Se7enSe7enSe7en/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot find JWT/bearer token", err)
		return
	}

	cfg.db.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token: token,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	respondWithJSON(w, http.StatusNoContent, nil)
}
