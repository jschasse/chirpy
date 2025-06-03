package main

import(
	"github.com/jschasse/chirpy/internal/auth"
	"log"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

func (cfg *apiConfig) handlerRefresh(writer http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting authorization header: %s", err)
		writer.WriteHeader(500)
		return
	}


	userID, err := cfg.dbQueries.GetUserFromRefreshToken(req.Context(), bearerToken)
	if err != nil {
		log.Printf("Refresh Token does not exist: %s", err)
		writer.WriteHeader(401)
		return
	}

	expiresIn := time.Hour

	accessToken, err := auth.MakeJWT(userID, cfg.secret, expiresIn)
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		writer.WriteHeader(500)
		return
	}

	type User struct{
		UserID uuid.UUID `json:"user_id"`
		Token  string `json:"token"`
	}

	user := User{
		UserID: userID,
		Token:	accessToken,
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling data: %s", err)
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(200)
	writer.Write(data)
}