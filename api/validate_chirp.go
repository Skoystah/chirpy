package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func validateChirp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, fmt.Sprintf(template, cfg.fileserverHits.Load()))
	if err != nil {
		log.Fatal(err)
	}
}
