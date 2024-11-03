package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func replaceProfanity(msg string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	wordSlice := strings.Split(msg, " ")
	newWordSlice := make([]string, len(wordSlice))
	// log.Printf("VIBE CHECK wordSlice: %v\n", wordSlice) // DEBUG

	for i, word := range wordSlice {
		_, ok := profaneWords[strings.ToLower(word)]
		if ok {
			newWordSlice[i] = "****"
			continue
		}
		newWordSlice[i] = word
	}

	// log.Printf("VIBE CHECK newWordSlice: %v\n", newWordSlice) // DEBUG
	return strings.Join(newWordSlice, " ")
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	maxCharLength := 140
	if len(params.Body) > maxCharLength {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}

	newMsg := replaceProfanity(params.Body)

	respondWithJSON(w, 200, struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: newMsg,
	})
}
