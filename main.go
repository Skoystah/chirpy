package main

import (
	"chirpy/internal/api"
	"chirpy/internal/config"
	"chirpy/internal/database"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	//Handle DB stuff
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading DB env: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	//Create API config
	//since its a very small struct a pointer is not needed - but it doesnt hurt to create it as a pointer
	//apiConfig := apiConfig{db: dbQueries}
	cfg := &config.ApiConfig{Db: dbQueries, Platform: platform, Secret: secret}

	//Handle HTTP server stuff
	serveMux := http.NewServeMux()
	server := &http.Server{Handler: serveMux, Addr: ":8080"}
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	//Register handlers
	serveMux.Handle("/app/", middlewareMetricsInc(cfg, fileServer))

	//GET
	serveMux.HandleFunc("GET /api/healthz", api.Healthz)
	serveMux.HandleFunc("GET /admin/metrics", api.Metrics(cfg))
	serveMux.HandleFunc("GET /api/chirps/{id}", api.GetChirp(cfg))
	serveMux.HandleFunc("GET /api/chirps", api.GetChirps(cfg))

	//POST
	serveMux.HandleFunc("POST /admin/reset", api.Reset(cfg))
	serveMux.HandleFunc("POST /api/login", api.Login(cfg))
	serveMux.HandleFunc("POST /api/users", api.CreateUser(cfg))
	serveMux.HandleFunc("POST /api/chirps", api.CreateChirp(cfg))
	serveMux.HandleFunc("POST /api/refresh", api.Refresh(cfg))
	serveMux.HandleFunc("POST /api/revoke", api.Revoke(cfg))

	//Start server
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func middlewareMetricsInc(cfg *config.ApiConfig, next http.Handler) http.Handler {
	//in this case we CONVERT a regular func(w, req) ... to a http.HandlerFunc TYPE. Its the same as converting int32 to int for example.
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits.Add(1)
		//we must call ServeHTTP manually in this case, otherwise the chain of handlers stops here.  e
		next.ServeHTTP(w, req)
	})
}
