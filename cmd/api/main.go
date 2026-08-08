package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"veterans-go-chi-server/internal/database"
	"veterans-go-chi-server/internal/handlers"
	"veterans-go-chi-server/internal/repositories"
	"veterans-go-chi-server/internal/services"
)

func main() {
	pool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	mediaRepository := repositories.NewMediaRepository(pool)
	mediaService := services.NewMediaService(mediaRepository)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bon dia desde Go Chi!")
	})
	
	r.Post("/media", mediaHandler.Upload)
	
	log.Println("Server is listening on port 8000...")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}