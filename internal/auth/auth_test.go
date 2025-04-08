package auth

import (
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
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

func TestJWT(t *testing.T) {
	userID := uuid.New()
	key := make([]byte, 32)
	rand.Read(key)
	tokenSecret := string(key)

	expiresIn, _ := time.ParseDuration("600s")

	signedString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("error creating JWT token: %v", err)
		t.Fail()
		return
	}

	userIDAfter, err := ValidateJWT(signedString, tokenSecret)
	if err != nil {
		t.Errorf("error decoding JWT token: %v", err)
		t.Fail()
		return
	}
	fmt.Println(userID, userIDAfter)
	if userIDAfter != userID {
		t.Errorf("error - userID before and after JWT not the same : %v - %v", userID, userIDAfter)
	}
}

func TestGetBearerToken(t *testing.T) {

	header1 := http.Header{}
	header1.Set("Authorization", "Bearer this is the tokenstring")
	header2 := http.Header{}

	tests := []struct {
		header    http.Header
		wantErr   bool
		wantToken string
	}{
		{
			header:    header1,
			wantErr:   false,
			wantToken: "this is the tokenstring",
		},
		{
			header:    header2,
			wantErr:   true,
			wantToken: "",
		},
	}

	for _, tt := range tests {
		tokenString, err := GetBearerToken(tt.header)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetBearerToken() error = %v, wantErr = %v", err, tt.wantErr)
		}
		if tokenString != tt.wantToken {
			t.Errorf("GetBearerToken() tokenString = %v, expected = %v", tokenString, tt.wantToken)
		}
	}
}
