package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"log"
	"net/http"
)

func CreateUser(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		params := model.CreateUserRequest{}

		err := decoder.Decode(&params)
		if err != nil {
			// an error will be thrown if the JSON is invalid or has the wrong types
			// any missing fields will simply have their values in the struct set to their zero value
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Printf("error creating user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newUser, err := db.CreateUserDB(cfg, model.User{Email: params.Email, HashedPassword: hashedPassword})
		if err != nil {
			log.Printf("error creating user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//you can also marshal but its more cumbersome for this purpose. Marshal is good when you need to save the
		//intermediate result.
		encoder := json.NewEncoder(w)

		response := model.CreateUserResponse{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		}
		w.WriteHeader(http.StatusCreated)
		err = encoder.Encode(&response)
		if err != nil {
			log.Printf("Error encoding parameters")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func UpdateUser(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// AUTH
		token, err := auth.GetBearerToken(req.Header)
		if err != nil {
			log.Printf("Error getting token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		authUserID, err := auth.ValidateJWT(token, cfg.Secret)
		if err != nil {
			log.Printf("Error validating token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(req.Body)
		params := model.UpdateUserRequest{}

		err = decoder.Decode(&params)
		if err != nil {
			// an error will be thrown if the JSON is invalid or has the wrong types
			// any missing fields will simply have their values in the struct set to their zero value
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Printf("error updating user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newUser, err := db.UpdateUserDB(cfg, model.User{ID: authUserID, Email: params.Email, HashedPassword: hashedPassword})
		if err != nil {
			log.Printf("error updating user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		encoder := json.NewEncoder(w)

		response := model.UpdateUserResponse{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		}
		w.WriteHeader(http.StatusOK)
		err = encoder.Encode(&response)
		if err != nil {
			log.Printf("Error encoding parameters")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
