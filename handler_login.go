package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/jschasse/chirpy/internal/auth"
	"github.com/jschasse/chirpy/internal/database"
	"github.com/google/uuid"
	"time"
)

func (cfg *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {
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

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error getting email: %s", err)
		writer.WriteHeader(401)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		log.Printf("Wrong password: %s", err)
		writer.WriteHeader(401)
		return
	}

	expiresIn := time.Hour

	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		writer.WriteHeader(500)
		return
	}

	refresh_token := auth.MakeRefreshToken()

	args := database.CreateRefreshTokenParams{
		Token: refresh_token,
		UserID: user.ID,
	}

	r_token, err := cfg.dbQueries.CreateRefreshToken(req.Context(), args)
	if err != nil {
		log.Printf("Error Creating Refresh Token in database: %s", err)
		writer.WriteHeader(500)
		return
	}


	type User struct {
		ID        	 uuid.UUID `json:"id"`
		CreatedAt 	 time.Time `json:"created_at"`
		UpdatedAt 	 time.Time `json:"updated_at"`
		Email     	 string    `json:"email"`
		Token	  	 string	   `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	resUser := User{
		ID:				user.ID,
		CreatedAt:		user.CreatedAt,
		UpdatedAt:		user.UpdatedAt,
		Email:			user.Email,
		Token:			token,
		RefreshToken:	r_token,
	}

	data, err := json.Marshal(resUser)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(200)
	writer.Write(data)

}