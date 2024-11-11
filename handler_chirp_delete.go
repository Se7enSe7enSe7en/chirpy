package main

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	var err error

	chirpId := r.PathValue("chirpId")
	chirpIdUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ChirpId is invalid UUID", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot get bearer token", err)
		return
	}
	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot validate access token", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpIdUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpIdUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
