package main

import(
	"net/http"
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"time"
	"github.com/jschasse/chirpy/internal/database"
	"sort"
)

func (cfg *apiConfig) handlerAllChirps(writer http.ResponseWriter, req *http.Request) {
	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

    authorIDStr := req.URL.Query().Get("author_id")
	sortStr := req.URL.Query().Get("sort")
    
	var resChirps []database.Chirp
	var err error

    if authorIDStr != "" {
        authorID, err := uuid.Parse(authorIDStr)
        if err != nil {
            log.Printf("Invalid author_id format: %s", err)
            writer.WriteHeader(400)
            return
        }
        
        resChirps, err = cfg.dbQueries.GetChirpsByUserID(req.Context(), authorID)
        if err != nil {
            log.Printf("Error getting chirps by user ID: %s", err)
            writer.WriteHeader(500)
            return
        }
    } else {
        resChirps, err = cfg.dbQueries.GetChirps(req.Context())
        if err != nil {
            log.Printf("Error getting chirps: %s", err)
            writer.WriteHeader(500)
            return
        }
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

	if sortStr == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
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