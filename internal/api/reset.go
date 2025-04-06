package api

import (
	"chirpy/internal/config"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Reset(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if cfg.Platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)

		_, err := io.WriteString(w, fmt.Sprint("Hits: ", cfg.FileserverHits.Load(), "\n"))
		if err != nil {
			log.Fatal(err)
		}

		err = cfg.Db.DeleteUsers(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		cfg.FileserverHits.Store(0)
	}
}
