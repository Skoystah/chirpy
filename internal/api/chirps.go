package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func CreateChirp(cfg *config.ApiConfig) http.HandlerFunc {
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
		params := model.CreateChirpRequest{}

		err = decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error decoding parameters")
			return
		}

		cleaned_chirp, err := validateChirp(params.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, err.Error())
			return
		}

		newChirp, err := db.CreateChirpDB(cfg, model.Chirp{Body: cleaned_chirp, UserID: authUserID})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error creating chirp")
			return
		}

		response := model.CreateChirpResponse{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
		}
		respondWithJSON(w, http.StatusCreated, response)
	})
}

func GetChirps(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		chirps, err := db.GetChirpsDB(cfg)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error fetching chirps")
			return
		}

		var response []model.GetChirpResponse

		for _, chirp := range chirps {
			response = append(response, model.GetChirpResponse{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}
		respondWithJSON(w, http.StatusOK, response)
	})
}

func GetChirp(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		chirpID, err := uuid.Parse(req.PathValue("id"))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error parsing ID")
			return
		}

		chirp, err := db.GetChirpDB(cfg, model.Chirp{ID: chirpID})
		if err != nil {
			respondWithError(w, http.StatusNotFound, err, "Chirp not found")
			return
		}

		response := model.GetChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		respondWithJSON(w, http.StatusOK, response)
	})
}
func DeleteChirp(cfg *config.ApiConfig) http.HandlerFunc {
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

		chirpID, err := uuid.Parse(req.PathValue("id"))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error parsing ID")
			return
		}

		chirp, err := db.GetChirpDB(cfg, model.Chirp{ID: chirpID})
		if err != nil {
			respondWithError(w, http.StatusNotFound, err, "Chirp not found")
			return
		}

		if chirp.UserID != authUserID {
			err = errors.New("Chirp does not belong to this user")
			respondWithError(w, http.StatusForbidden, err, err.Error())
			return
		}

		err = db.DeleteChirpDB(cfg, chirp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Error deleting chirp")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
func validateChirp(chirp string) (string, error) {

	const maxChirpLength = 140
	if chirpLength := len(chirp); chirpLength > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	//gdo: I first had a map[string]bool with 'true' for all the values. This also works but by convention struct{} is more used,
	//also empty structs do not take up any memory. Bool would be better if you want active/non-active values in the list for example.
	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(chirp, " ")
	for i, word := range words {
		if _, ok := profanity[strings.ToLower(word)]; ok {
			words[i] = strings.Repeat("*", 4)
		}
	}
	return strings.Join(words, " "), nil
}
