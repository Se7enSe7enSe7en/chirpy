package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpList(w http.ResponseWriter, r *http.Request) {
	chirpList, err := cfg.db.GetChirpList(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp list", err)
		return
	}

	responseList := make([]Chirp, len(chirpList))
	for i, chirp := range chirpList {
		responseList[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, responseList)
}
