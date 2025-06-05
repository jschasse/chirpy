package main

import(
	"github.com/google/uuid"
	"net/http"
	"encoding/json"
	"log"
	"github.com/jschasse/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks(writer http.ResponseWriter, req *http.Request) {
	type reqData struct {
		UserID string `json:"user_id"`
	}

	type body struct {
		Event string `json:"event"`
		Data reqData `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Printf("Error getting authorization header: %s", err)
		writer.WriteHeader(401)
		return
	}

	if apiKey != cfg.apiKey {
		log.Printf("APIKey does not match")
		writer.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := body{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
	}

	if params.Event != "user.upgraded" {
		log.Printf("Dont care about this event")
		writer.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
    	log.Printf("Invalid user ID format: %s", err)
    	writer.WriteHeader(400)
    	return
	}

	err = cfg.dbQueries.UpgradeUserRed(req.Context(), userID)
	if err != nil {
		log.Printf("User cant be found: %s", err)
		writer.WriteHeader(404)
		return
	}

	writer.WriteHeader(204)
}