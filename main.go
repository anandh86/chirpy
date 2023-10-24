package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	const filepathRoot = "."
	const port = "8080" // Set your desired port

	db, err := NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	r := chi.NewRouter()

	r.Use(middlewareCors)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)

	subRouter := chi.NewRouter()

	subRouter.Get("/healthz", handlerReadiness)
	subRouter.Get("/reset", apiCfg.handlerReset)
	subRouter.Post("/chirps", apiCfg.postChirp)
	subRouter.Get("/chirps", apiCfg.getChirp)
	subRouter.Get("/chirps/{chirpID}", apiCfg.getChirpById)

	r.Mount("/admin", adminRouter)
	r.Mount("/api", subRouter)

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
