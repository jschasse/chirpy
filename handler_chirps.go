package main 

import(
	"net/http"
	"encoding/json"
	"log"
	"strings"
	"github.com/jschasse/chirpy/internal/database"
	"github.com/jschasse/chirpy/internal/auth"
	"github.com/google/uuid"
	"time"
)

func (cfg *apiConfig) handlerChirps(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	type reqChirp struct {
		Body string `json:"body"`
	}

	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	type errorVals struct {
		Error string `json:"error"`
	}

	type cleanedBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqChirp{}
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

	params.Body = joinedSentence

	bearer, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting Bearer: %s", err)
		writer.WriteHeader(500)
		return
	}

	userID, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		log.Printf("Invalid JWT: %s", err)
		writer.WriteHeader(401)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:	params.Body,
		UserID:	userID,
	}

	resChirp, err := cfg.dbQueries.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		writer.WriteHeader(500)
		return
	}

	mainChirp := chirp{
		ID:			resChirp.ID,
		CreatedAt:	resChirp.CreatedAt,
		UpdatedAt:	resChirp.UpdatedAt,
		Body:		resChirp.Body,
		UserID:		resChirp.UserID,
	}

	data, err := json.Marshal(mainChirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(201)
	writer.Write(data)

}