package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
	"github.com/Se7enSe7enSe7en/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	accessToken, err := auth.MakeJWT(
		userDB.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot make JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot make Refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userDB.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60), // 60 days
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot create refresh token", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          userDB.ID,
			CreatedAt:   userDB.CreatedAt,
			UpdatedAt:   userDB.UpdatedAt,
			Email:       userDB.Email,
			IsChirpyRed: userDB.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
