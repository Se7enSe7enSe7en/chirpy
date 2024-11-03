package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	isDev := false
	if os.Getenv("PLATFORM") == "dev" {
		isDev = true
	}

	if !isDev {
		w.WriteHeader(403)
		return
	}

	// Reset user table
	err := cfg.db.ResetUserTable(r.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
	}

	// Reset fileserverHits state
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

// Test: curl -X POST http://localhost:8080/admin/reset
