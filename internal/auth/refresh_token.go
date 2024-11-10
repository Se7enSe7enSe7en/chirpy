package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	// generate a random 256-bit (32 bytes) hex-encoded string
	rand32bytes := make([]byte, 32)
	_, err := rand.Read(rand32bytes)
	if err != nil {
		return "", err
	}
	rand32bytesStr := fmt.Sprintf("%v", rand32bytes)

	return hex.EncodeToString([]byte(rand32bytesStr)), nil // TMP
}
