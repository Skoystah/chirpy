package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func Refresh(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokenString, err := auth.GetBearerToken(req.Header)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error reading token")
			return
		}

		refreshToken, err := db.GetRefreshTokenDB(cfg, model.RefreshToken{Token: tokenString})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err, "Error retrieving token")
			return
		}

		if refreshToken.Expires_at.UTC().Before(time.Now().UTC()) {
			err = errors.New("Token has expired")
			respondWithError(w, http.StatusUnauthorized, err, err.Error())
			return
		}

		fmt.Println(refreshToken.Revoked_at)
		if refreshToken.Revoked_at.Valid {
			err = errors.New("Token has been revoked")
			respondWithError(w, http.StatusUnauthorized, err, err.Error())
			return
		}

		const expiresIn = time.Hour
		token, err := auth.MakeJWT(refreshToken.UserID, cfg.Secret, expiresIn)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error creating JWT token")
			return
		}

		response := model.RefreshResponse{
			Token: token,
		}
		respondWithJSON(w, http.StatusOK, response)
	})
}

func Revoke(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokenString, err := auth.GetBearerToken(req.Header)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error retrieving token")
			return
		}

		err = db.RevokeRefreshTokenDB(cfg, model.RefreshToken{Token: tokenString})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error revoking token")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
