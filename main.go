package main

import (
	"log"
	"net/http"
)

type apiHandler struct {
}

func (a apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	log.Fatal("unimplemented")
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// create server mux
	mux := http.NewServeMux()
	// handle root path "/" then use http.FileServer to call the index.html in the "." directory
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	// init the server config
	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start the server
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatalln(s.ListenAndServe())
}
