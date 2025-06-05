package main

import(
	"net/http"
	"encoding/json"
	"log"
	"github.com/jschasse/chirpy/internal/auth"
	"github.com/jschasse/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting Bearer: %s", err)
		writer.WriteHeader(401)
		return
	}

	_, err = auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		log.Printf("Invalid JWT: %s", err)
		writer.WriteHeader(401)
		return
	}

	type reqEmail struct {
		Email string	`json:"email"`
		Password string	`json:"password"`
	}

	type resEmail struct {
		Email string	`json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqEmail{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		writer.WriteHeader(500)
		return
	}

	p := database.InsertPasswordParams{
		HashedPassword: hashedPassword,
		Email:			params.Email,
	}

	err = cfg.dbQueries.InsertPassword(req.Context(), p)
	if err != nil {
		log.Printf("Error inserting password: %s", err)
		writer.WriteHeader(500)
		return
	}

	resP := resEmail{
		Email: params.Email,
	}

	data, err := json.Marshal(resP)
	if err != nil {
		log.Printf("Error Marshaling data: %s", err)
		writer.WriteHeader(500)
		return
	}


	writer.WriteHeader(200)
	writer.Write(data)

}