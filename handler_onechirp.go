package main

import(
	"net/http"
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"time"
)

func (cfg *apiConfig) handlerOneChirp(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirpIDstr := req.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
        http.Error(writer, "Invalid chirp ID", http.StatusBadRequest)
        return
    }

	chirpdata, err := cfg.dbQueries.GetOneChirp(req.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		writer.WriteHeader(500)
		return
	}

	tmp := chirp{
			ID:			chirpdata.ID,
			CreatedAt:	chirpdata.CreatedAt,
			UpdatedAt:	chirpdata.UpdatedAt,
			Body:		chirpdata.Body,
			UserID:		chirpdata.UserID,
		}

	data, err := json.Marshal(tmp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(200)
    writer.Write(data)
}