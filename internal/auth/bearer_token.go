package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if bearerToken == "" {
		return "", errors.New("no bearer token from headers")
	}

	return strings.ReplaceAll(bearerToken, "Bearer ", ""), nil
}
