package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"net/http"
)

func CreateUser(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		params := model.CreateUserRequest{}

		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error decoding parameters")
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error hashing password")
			return
		}

		newUser, err := db.CreateUserDB(cfg, model.User{Email: params.Email, HashedPassword: hashedPassword})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error creating user")
			return
		}

		response := model.CreateUserResponse{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		}
		respondWithJSON(w, http.StatusCreated, response)
	})
}

func UpdateUser(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// AUTH
		token, err := auth.GetBearerToken(req.Header)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error reading JWT token")
			return
		}
		authUserID, err := auth.ValidateJWT(token, cfg.Secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err, "Error validating JWT token")
			return
		}

		decoder := json.NewDecoder(req.Body)
		params := model.UpdateUserRequest{}

		err = decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error decoding parameters")
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error hashing password")
			return
		}

		newUser, err := db.UpdateUserDB(cfg, model.User{ID: authUserID, Email: params.Email, HashedPassword: hashedPassword})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error updating user")
			return
		}

		response := model.UpdateUserResponse{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		}
		respondWithJSON(w, http.StatusOK, response)
	})
}
