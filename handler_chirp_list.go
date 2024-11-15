package main

import (
	"net/http"
	"slices"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpList(w http.ResponseWriter, r *http.Request) {
	var err error

	authorId := uuid.Nil
	authorIdStr := r.URL.Query().Get("author_id")
	if authorIdStr != "" {
		authorId, err = uuid.Parse(authorIdStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Cannot parse author_id query param", err)
		}
	}

	type sortOrderType string
	const (
		sortOrderAsc  sortOrderType = "asc"
		sortOrderDesc sortOrderType = "desc"
	)
	var sortOrder sortOrderType = sortOrderAsc
	sortOrderStr := r.URL.Query().Get("sort")
	if sortOrderStr == string(sortOrderDesc) {
		sortOrder = sortOrderDesc
	}

	chirpList, err := cfg.db.GetChirpList(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp list", err)
		return
	}

	responseList := make([]Chirp, len(chirpList))
	for i, chirp := range chirpList {
		// optional filter: authorId
		if authorId != uuid.Nil && authorId != chirp.UserID {
			// if authorId exists and authorId is not equal to the chirp.UserID
			// then skip
			continue
		}

		responseList[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	if sortOrder == sortOrderDesc {
		slices.Reverse(responseList)
	}

	respondWithJSON(w, http.StatusOK, responseList)
}
