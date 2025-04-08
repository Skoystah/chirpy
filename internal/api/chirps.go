package api

import (
	"chirpy/internal/config"
	"chirpy/internal/db"
	"chirpy/internal/model"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

func CreateChirp(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		decoder := json.NewDecoder(req.Body)
		params := model.CreateChirpRequest{}

		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		valid, cleaned_chirp := validateChirp(params.Body)
		if !valid {
			//Todo proper way to send error out
			w.WriteHeader(http.StatusInternalServerError)
		}

		newChirp, err := db.CreateChirpDB(cfg, model.Chirp{Body: cleaned_chirp, UserID: params.UserID})
		if err != nil {
			log.Printf("error creating chirp: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//you can also marshal but its more cumbersome for this purpose. Marshal is good when you need to save the
		//intermediate result.
		encoder := json.NewEncoder(w)

		response := model.CreateChirpResponse{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
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

func GetChirps(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		chirps, err := db.GetChirpsDB(cfg)
		if err != nil {
			log.Printf("error fetching chirps: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//you can also marshal but its more cumbersome for this purpose. Marshal is good when you need to save the
		//intermediate result.
		encoder := json.NewEncoder(w)

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

		w.WriteHeader(http.StatusOK)
		err = encoder.Encode(&response)
		if err != nil {
			log.Printf("Error encoding parameters")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func GetChirp(cfg *config.ApiConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		chirpID, err := uuid.Parse(req.PathValue("id"))
		if err != nil {
			log.Printf("error parsing uuid to string: %v", err)
		}

		chirp, err := db.GetChirpDB(cfg, model.Chirp{ID: chirpID})
		if err != nil {
			log.Printf("error fetching chirp: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//you can also marshal but its more cumbersome for this purpose. Marshal is good when you need to save the
		//intermediate result.
		encoder := json.NewEncoder(w)

		response := model.GetChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
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
func validateChirp(chirp string) (bool, string) {

	const maxChirpLength = 140
	if chirpLength := len(chirp); chirpLength > maxChirpLength {
		log.Printf("Chirp is too long: %d - %s", chirpLength, chirp)
		return false, chirp
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
	return true, strings.Join(words, " ")
}
