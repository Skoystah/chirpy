package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Login(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		params := model.LoginRequest{}

		err := decoder.Decode(&params)
		if err != nil {
			// an error will be thrown if the JSON is invalid or has the wrong types
			// any missing fields will simply have their values in the struct set to their zero value
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := db.GetUserDB(cfg, model.User{Email: params.Email})
		if err != nil {
			log.Printf("error creating user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("Password", params.Password, "Hash", user.HashedPassword)
		err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
		if err != nil {
			log.Print("password not correct")
			w.WriteHeader(http.StatusUnauthorized)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")

			_, err := io.WriteString(w, "Incorrect email or password")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		//you can also marshal but its more cumbersome for this purpose. Marshal is good when you need to save the
		//intermediate result.
		encoder := json.NewEncoder(w)

		response := model.CreateUserResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
		err = encoder.Encode(&response)
		if err != nil {
			log.Printf("Error encoding parameters")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}
