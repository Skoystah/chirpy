package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"errors"
	"net/http"
)

func PolkaWebhooks(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// AUTH
		apiKey, err := auth.GetAPIKey(req.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err, "Error reading API Key")
			return
		}

		if apiKey != cfg.PolkaKey {
			err = errors.New("API Key for Polka not correct")
			respondWithError(w, http.StatusUnauthorized, err, err.Error())
			return
		}

		decoder := json.NewDecoder(req.Body)
		params := model.PolkaWebhookRequest{}

		err = decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error decoding parameters")
			return
		}

		switch params.Event {
		case "user.upgraded":
			_, err := db.UpdateUserRedDB(cfg, model.User{ID: params.Data.UserID, IsChirpyRed: true})
			if err != nil {
				respondWithError(w, http.StatusNotFound, err, "Error updating user")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNoContent)
			return
		}

	})
}
