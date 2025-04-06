package api

import (
	"io"
	"log"
	"net/http"
)

func Healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, "OK")
	if err != nil {
		log.Fatal(err)
	}
}
