package api

import (
	"chirpy/internal/model"
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&data)
	if err != nil {
		log.Printf("Error encoding parameters: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, err error, msg string) {
	if err != nil {
		log.Println(err)
	}
	if statusCode > 499 {
		log.Printf("Responding with error 5XX: %s", msg)
	}
	respondWithJSON(w, statusCode, model.ErrorResponse{Error: msg})
}
