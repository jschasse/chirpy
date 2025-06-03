package main

import(
	"github.com/jschasse/chirpy/internal/auth"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(writer http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting authorization header: %s", err)
		writer.WriteHeader(500)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), bearerToken)
	if err != nil {
		log.Printf("Error revoking token: %s", err)
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(204)
}