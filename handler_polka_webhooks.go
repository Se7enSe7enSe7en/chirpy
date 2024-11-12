package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserId string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decoder.Decode(&params)

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot parse user id", err)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cannot find user", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
