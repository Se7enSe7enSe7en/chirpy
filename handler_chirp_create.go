package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
	"github.com/Se7enSe7enSe7en/chirpy/internal/database"
	"github.com/google/uuid"
)

func getCleanedBody(body string, badWords map[string]struct{}) string {
	wordSlice := strings.Split(body, " ")

	for i, word := range wordSlice {
		loweredWord := strings.ToLower(word)
		_, ok := badWords[loweredWord]
		if ok {
			wordSlice[i] = "****"
			continue
		}
		wordSlice[i] = word
	}

	cleanedBody := strings.Join(wordSlice, " ")
	return cleanedBody
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		ID        uuid.UUID     `json:"id"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
		Body      string        `json:"body"`
		UserID    uuid.NullUUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot get bearer token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT invalid or expired", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	maxCharLength := 140
	if len(params.Body) > maxCharLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedBody := getCleanedBody(params.Body, badWords)

	chirpDB, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: uuid.NullUUID{UUID: userId, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		ID:        chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body:      chirpDB.Body,
		UserID:    chirpDB.UserID,
	})
}
