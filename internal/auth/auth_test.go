package auth

import (
	"testing"
)

func TestAuth(t *testing.T) {
	password := "weakpassword"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("error hashing password: %v - %v", password, err)
		t.Fail()
		return
	}

	err = CheckPasswordHash(hashedPassword, password)
	if err != nil {
		t.Errorf("error comparing hash %v and password %v - %v", hashedPassword, password, err)
		t.Fail()
	}
}
