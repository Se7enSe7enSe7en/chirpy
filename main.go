package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Se7enSe7enSe7en/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// initialize DB
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("Error in opening DB")
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
	}
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	// create server mux
	mux := http.NewServeMux()

	// register handlers
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// init the server config
	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start the server
	log.Printf("Serving on port: %s\n", port)
	// log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatalln(s.ListenAndServe())
}
