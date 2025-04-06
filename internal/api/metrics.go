package api

import (
	"chirpy/internal/config"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Metrics(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		template := `<html>
	<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
	</body>
	</html>`

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_, err := io.WriteString(w, fmt.Sprintf(template, cfg.FileserverHits.Load()))
		if err != nil {
			log.Fatal(err)
		}
	}
}
