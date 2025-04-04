package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (cfg *apiConfig) metrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, fmt.Sprint("Hits: ", cfg.fileserverHits.Load(), "\n"))
	if err != nil {
		log.Fatal(err)
	}
}
