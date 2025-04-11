package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"net/http"
	"time"
)

func Login(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		params := model.LoginRequest{}

		err := decoder.Decode(&params)
		if err != nil {
			// an error will be thrown if the JSON is invalid or has the wrong types
			// any missing fields will simply have their values in the struct set to their zero value
			respondWithError(w, http.StatusInternalServerError, err, "Error reading request")
			return
		}

		user, err := db.GetUserDB(cfg, model.User{Email: params.Email})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error creating user")
			return
		}

		err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err, "Incorrect email or password")
			return
		}

		const expiresIn = time.Hour
		token, err := auth.MakeJWT(user.ID, cfg.Secret, expiresIn)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error creating JWT token")
			return
		}

		//Create REFRESH token
		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error creating refresh token")
			return
		}

		const refreshExpiresIn = 60
		err = db.CreateRefreshToken(cfg, model.RefreshToken{
			Token:      refreshToken,
			UserID:     user.ID,
			Expires_at: time.Now().AddDate(0, 0, 60),
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error storing refresh token")
			return
		}

		response := model.LoginResponse{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshToken,
		}
		respondWithJSON(w, http.StatusOK, response)
	})
}
