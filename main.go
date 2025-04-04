package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()

	server := &http.Server{Handler: serveMux, Addr: ":8080"}

	//new handler = File server - path = "." or current directory
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	//add handler to server multiplexer
	serveMux.Handle("/app/", fileServer)
	serveMux.HandleFunc("/healthz", healthz)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, "OK")
	if err != nil {
		log.Fatal(err)
	}

}
