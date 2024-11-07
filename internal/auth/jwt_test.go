package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	id1 := uuid.New()
	key := "4zXPK6Vzm2mcf3tAFqxnBhq/yqfnGauBhr0BgkMBn9DeKtacBH88nYvvXJfXjn4egJnfIJjqj5mhDW3zhbIsaw=="
	expiry := time.Second

	jwtToken, err := MakeJWT(id1, key, expiry)
	if err != nil {
		t.Fatalf("cannot make jwt: %v", err)
	}

	id2, err := ValidateJWT(jwtToken, key)
	if err != nil {
		t.Fatalf("cannot validate jwt: %v", err)
	}

	// test if we got the same id
	if id1 != id2 {
		t.Errorf("id1 and id2 should be equal, id1 = %v, id2 = %v", id1, id2)
	}

	// test if token will expire
	time.Sleep(expiry)
	_, err = ValidateJWT(jwtToken, key)
	if err == nil {
		t.Errorf("token did not expire")
	}
}
