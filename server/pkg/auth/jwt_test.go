package auth

import "testing"

func TestJwt(t *testing.T) {
	jwt, err := GenerateJWT("admin")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(jwt)
	claims, err := ParseJWT(jwt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
