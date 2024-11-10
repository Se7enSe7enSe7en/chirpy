package auth

import "testing"

func TestMakeRefreshToken(t *testing.T) {
	refreshToken, err := MakeRefreshToken()
	if err != nil {
		t.Errorf("make refresh token error:  %v", err)
	}
	t.Logf("%v", refreshToken)
}
