package auth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	token1 := "4zXPK6Vzm2mcf3tAFqxnBhq/yqfnGauBhr0BgkMBn9DeKtacBH88nYvvXJfXjn4egJnfIJjqj5mhDW3zhbIsaw=="

	header := http.Header{
		"Authorization": {fmt.Sprintf("Bearer %v", token1)},
	}

	token2, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("Cannot get bearer token from header: %v", err)
	}

	// test if token1 and token2 is the same
	if token1 != token2 {
		t.Errorf("token1 and token2 is not the same, token1 = %v, token2 = %v", token1, token2)
	}
}
