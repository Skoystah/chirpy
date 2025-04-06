package api

import (
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

		newUser, err := db.CreateUserDB(cfg, model.User{Email: params.Email})
		if err != nil {
			log.Printf("error creating user: %w", err)
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
