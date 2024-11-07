package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Use jwt.NewWithClaims to create a new token.
	// Use jwt.SigningMethodHS256 as the signing method.
	// Use jwt.RegisteredClaims as the claims.
	// Set the Issuer to "chirpy"
	// Set IssuedAt to the current time in UTC
	// Set ExpiresAt to the current time plus the expiration time (expiresIn)
	// Set the Subject to a stringified version of the user's id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprint(userId),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, errors.New("token is invalid or has expired")
	}

	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, errors.New("cannot get subject from token")
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.UUID{}, errors.New("cannot parse the userId from token")
	}

	return userId, nil
}
