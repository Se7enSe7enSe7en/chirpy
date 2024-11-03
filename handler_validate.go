package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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

	respondWithJSON(w, http.StatusOK, response{
		CleanedBody: cleanedBody,
	})
}
