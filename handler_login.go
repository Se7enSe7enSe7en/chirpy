package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)

	userDB, err := cfg.db.LoginUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, userDB.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	jwtToken, err := auth.MakeJWT(
		userDB.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot make JWT token", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        userDB.ID,
			CreatedAt: userDB.CreatedAt,
			UpdatedAt: userDB.UpdatedAt,
			Email:     userDB.Email,
		},
		Token: jwtToken,
	})
}
