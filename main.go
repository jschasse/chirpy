package main

import(
	"fmt"
	"net/http"
	"sync/atomic"
)


type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const port = "8080"
	apiCfg := &apiConfig{}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /healthz", handlerHealth)
	serveMux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /reset", apiCfg.handlerReset)
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	
	

	server := http.Server{
		Addr:		":" + port,
		Handler:	serveMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("%s", err)
	}


}

func handlerHealth(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	hits := cfg.fileServerHits.Load()
	writer.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reset"))
}