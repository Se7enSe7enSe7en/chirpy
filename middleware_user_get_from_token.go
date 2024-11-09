package main

import (
	"context"
	"net/http"

	"github.com/Se7enSe7enSe7en/chirpy/internal/auth"
	"github.com/Se7enSe7enSe7en/chirpy/internal/constants"
)

func (cfg *apiConfig) middlewareGetUserFromToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Cannot get bearer token", err)
			return
		}

		userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "JWT invalid or expired", err)
			return
		}

		// create userId variable in r.Context()
		ctx := context.WithValue(r.Context(), constants.ContextVariableUserId{}, userId)
		// note: using a struct as the context key is ok: https://stackoverflow.com/a/68100270

		// pass the context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
