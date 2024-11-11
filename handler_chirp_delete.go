package main

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
	"github.com/Se7enSe7enSe7en/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	var err error

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Cannot get bearer token", err)
		return
	}
	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 403, "Cannot validate access token", err)
		return
	}

	chirpId := r.PathValue("chirpId")
	chirpIdUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ChirpId is invalid UUID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpIdUUID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, 403, "Chirp doesn't belong to the user", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		UserID: userId,
		ID:     chirpIdUUID,
	})
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
