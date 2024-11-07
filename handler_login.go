package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds string `json:"expires_in_seconds,omitempty"`
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

	expiryStr := params.ExpiresInSeconds
	if expiryStr == "" {
		expiryStr = "3600" // 1hr in seconds
	}

	expiry, err := time.ParseDuration(fmt.Sprintf("%vs", expiryStr))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot parse expires_in_seconds", err)
	}

	jwtToken, err := auth.MakeJWT(userDB.ID, cfg.token, expiry)
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
