package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"io"
	"log"
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
		//TODO MORE ELEGANTLY!
		const expiresIn = time.Hour
		token, err := auth.MakeJWT(user.ID, cfg.Secret, expiresIn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}

		//Create REFRESH token
		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		const refreshExpiresIn = 60
		err = db.CreateRefreshToken(cfg, model.RefreshToken{
			Token:      refreshToken,
			UserID:     user.ID,
			Expires_at: time.Now().AddDate(0, 0, 60),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		//you can also marshal but its more cumbersome for this purpose. Marshal is good when you need to save the
		//intermediate result.
		encoder := json.NewEncoder(w)

		response := model.LoginResponse{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshToken,
		}
		err = encoder.Encode(&response)
		if err != nil {
			log.Printf("Error encoding parameters")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}
