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

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

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
	var err error

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
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
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		UserID:    chirpDB.UserID,
		Body:      chirpDB.Body,
	})
}
