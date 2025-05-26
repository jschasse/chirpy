package main

import _ "github.com/lib/pq"

import(
	"fmt"
	"net/http"
	"sync/atomic"
	"encoding/json"
	"log"
	"strings"
	"os"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/jschasse/chirpy/internal/database"
	"github.com/google/uuid"
	"time"
)


type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries *database.Queries
}

func main() {
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Print(err)
	}

	apiCfg := &apiConfig{}

	apiCfg.dbQueries = database.New(db)

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /api/healthz", handlerHealth)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
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

func handlerHealth(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(writer, r)
	})
}

func (cfg *apiConfig) handlerMetrics(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	htmlTemplate := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	hits := cfg.fileServerHits.Load()
	writer.Write([]byte(fmt.Sprintf(htmlTemplate, hits)))
}

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("Error deleting Users: %s", err)
		writer.WriteHeader(500)
		return
	}
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reset and Users deleted"))
}

func handlerValidate(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	type parameters struct {
		Body string `json:"body"`
	}

	type cleanedBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	type errorVals struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
	}

	

	if len(params.Body) > 140 {
		invalid := errorVals{
			Error: "Chirp is too long",
		}

		dat, err := json.Marshal(invalid)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			writer.WriteHeader(500)
			return
		}
		writer.WriteHeader(400)
    	writer.Write(dat)
		return
	}

	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	
	splitSentence := strings.Split(params.Body, " ")

	for i := 0; i < len(splitSentence); i++ {
		for j := 0; j < len(bannedWords); j++ {
			if strings.Contains(strings.ToLower(splitSentence[i]), bannedWords[j]) {
				splitSentence[i] = strings.ReplaceAll(strings.ToLower(splitSentence[i]), bannedWords[j], "****")
			}
		}
	}

	joinedSentence := strings.Join(splitSentence, " ")

	body := cleanedBody{
		CleanedBody: joinedSentence,
	}

	dat, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(500)
		return
	}
	writer.WriteHeader(200)
    writer.Write(dat)
	

}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	type email struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := email{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
	}

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	databseUser, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		writer.WriteHeader(500)
		return
	}

	user := User{
		ID:			databseUser.ID,
		CreatedAt:	databseUser.CreatedAt,
		UpdatedAt:	databseUser.UpdatedAt,
		Email:		databseUser.Email,
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(500)
		return
	}


	writer.WriteHeader(201)
	writer.Write(data)
}
