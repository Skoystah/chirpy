package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	serveMux := http.NewServeMux()
	//since its a very small struct a pointer is not needed - but it doesnt hurt to create it as a pointer
	apiConfig := apiConfig{}

	server := &http.Server{Handler: serveMux, Addr: ":8080"}

	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", apiConfig.middlewareMetricsInc(fileServer))

	serveMux.HandleFunc("GET /healthz", healthz)
	serveMux.HandleFunc("GET /metrics", apiConfig.metrics)
	serveMux.HandleFunc("POST /reset", apiConfig.reset)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	//in this case we CONVERT a regular func(w, req) ... to a http.HandlerFunc TYPE. Its the same as converting int32 to int for example.
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		//we must call ServeHTTP manually in this case, otherwise the chain of handlers stops here.  e
		next.ServeHTTP(w, req)
	})
}
