package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	// create server mux
	mux := http.NewServeMux()

	// init the server
	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start the server
	log.Printf("Serving in port: %v\n", port)
	log.Fatalln(s.ListenAndServe())
}
