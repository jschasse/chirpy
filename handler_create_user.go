package main

import(
	"net/http"
	"encoding/json"
	"log"
	"github.com/jschasse/chirpy/internal/database"
	"github.com/jschasse/chirpy/internal/auth"
	"github.com/google/uuid"
	"time"
)

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	type email struct {
		Email string `json:"email"`
		Password string `json:"password"`
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		writer.WriteHeader(500)
		return
	}

	databseUser, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		writer.WriteHeader(500)
		return
	}


	insertParams := database.InsertPasswordParams{
		HashedPassword: hashedPassword,
		Email:			params.Email,
	}

	err = cfg.dbQueries.InsertPassword(req.Context(), insertParams)
	if err != nil {
		log.Printf("Error Inserting password: %s", err)
		writer.WriteHeader(500)
		return
	}

	user := User{
		ID:				databseUser.ID,
		CreatedAt:		databseUser.CreatedAt,
		UpdatedAt:		databseUser.UpdatedAt,
		Email:			databseUser.Email,
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