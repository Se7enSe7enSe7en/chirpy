package main

import (
	"net/http"
	"time"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot find JWT/bearer token", err)
		return
	}

	dbToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Cannot find token in database", err)
	}

	if time.Now().Compare(dbToken.ExpiresAt) == 1 || dbToken.RevokedAt.Valid {
		respondWithError(w, 401, "Token is expired", nil)
	}

	userDB, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot get user from refresh token in the database", err)
	}

	newToken, err := auth.MakeJWT(userDB.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot make JWT", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}
