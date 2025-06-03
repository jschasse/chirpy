package main

import _ "github.com/lib/pq"

import(
	"fmt"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/jschasse/chirpy/internal/database"
)


type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries *database.Queries
	platform string
	secret string
}

func main() {
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretString := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Print(err)
	}

	apiCfg := &apiConfig{}

	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = platform
	apiCfg.secret = secretString

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /api/healthz", handlerHealth)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerOneChirp)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	
	

	server := http.Server{
		Addr:		":" + port,
		Handler:	serveMux,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("%s", err)
	}


}


















