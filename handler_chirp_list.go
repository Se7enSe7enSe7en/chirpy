package main

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpList(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Body []database.Chirp `json:"body"`
	}

	chirpList, err := cfg.db.GetChirpList(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp list", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Body: chirpList,
	})
}
