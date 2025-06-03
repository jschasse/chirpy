package main

import(
	"net/http"
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"time"
)

func (cfg *apiConfig) handlerAllChirps(writer http.ResponseWriter, req *http.Request) {
	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	resChirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		writer.WriteHeader(500)
		return
	}

	chirps := []chirp{}

	for i := 0; i < len(resChirps); i++ {
		tmp := chirp{
			ID:			resChirps[i].ID,
			CreatedAt:	resChirps[i].CreatedAt,
			UpdatedAt:	resChirps[i].UpdatedAt,
			Body:		resChirps[i].Body,
			UserID:		resChirps[i].UserID,
		}

		chirps = append(chirps, tmp)
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(200)
    writer.Write(data)
}