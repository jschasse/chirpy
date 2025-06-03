package main

import(
	"net/http"
	"log"
)

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev"{
		log.Printf("Not in local enviornment")
		writer.WriteHeader(403)
		return
	}
	cfg.fileServerHits.Store(0)
	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("Error deleting Users: %s", err)
		writer.WriteHeader(500)
		return
	}
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reset and Users deleted"))
}