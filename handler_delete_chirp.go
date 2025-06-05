package main

import(
	"net/http"
	"log"
	"github.com/jschasse/chirpy/internal/auth"
	"github.com/jschasse/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(writer http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting Bearer: %s", err)
		writer.WriteHeader(401)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		log.Printf("Invalid JWT: %s", err)
		writer.WriteHeader(401)
		return
	}

	chirpIDstr := req.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
        http.Error(writer, "Invalid chirp ID", http.StatusBadRequest)
        return
    }

	chirp, err := cfg.dbQueries.GetOneChirp(req.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		writer.WriteHeader(404)
		return
	}

	if chirp.UserID != userID {
		log.Printf("User is not author of chirp: %s", err)
		writer.WriteHeader(403)
		return
	}

	params := database.DeleteChirpParams{
		ID: chirpID,
		UserID: userID,
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), params)
	if err != nil {
		log.Printf("Invalid Chirp: %s", err)
		writer.WriteHeader(404)
		return
	}

	writer.WriteHeader(204)

}