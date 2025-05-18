package main

import(
	"fmt"
	"net/http"
)


func main() {
	const port = "8080"


	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Addr:		":" + port,
		Handler:	serveMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("%s", err)
	}


}