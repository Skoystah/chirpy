package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Refresh(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokenString, err := auth.GetBearerToken(req.Header)
		if err != nil {
			log.Printf("error retrieving token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		refreshToken, err := db.GetRefreshTokenDB(cfg, model.RefreshToken{Token: tokenString})
		if err != nil {
			log.Printf("error retrieving token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//if refreshToken.Expires_at.Compare(time.Now()) <= 0 {
		if refreshToken.Expires_at.UTC().Before(time.Now().UTC()) {
			log.Printf("Token has expired: %v", refreshToken.Expires_at)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println(refreshToken.Revoked_at)
		if refreshToken.Revoked_at.Valid {
			log.Printf("Token has been revoked at %v", refreshToken.Revoked_at.Time)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		const expiresIn = time.Hour
		token, err := auth.MakeJWT(refreshToken.UserID, cfg.Secret, expiresIn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		response := model.RefreshResponse{
			Token: token,
		}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(&response)
		if err != nil {
			log.Printf("Error encoding parameters")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func Revoke(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokenString, err := auth.GetBearerToken(req.Header)
		if err != nil {
			log.Printf("error retrieving token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = db.RevokeRefreshTokenDB(cfg, model.RefreshToken{Token: tokenString})
		if err != nil {
			log.Printf("error revoking token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
