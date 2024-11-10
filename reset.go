package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	// Reset user table
	err := cfg.db.ResetUserTable(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset the user table", err)
	}

	// // Reset chirps table
	// err = cfg.db.ResetChirpsTable(r.Context())
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't reset the chirps table", err)
	// }

	// // Reset refresh_tokens table
	// err = cfg.db.ResetRefreshTokensTable(r.Context())
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't reset the refresh tokens table", err)
	// }

	// Reset fileserverHits state
	cfg.fileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}

// Test: curl -X POST http://localhost:8080/admin/reset
