package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"veterans-go-chi-server/internal/handlers"
	"veterans-go-chi-server/internal/services"
)

func main() {
	r := chi.NewRouter()

	mediaService := services.NewMediaService()

	mediaHandler := handlers.NewMediaHandler(mediaService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bon dia desde Go Chi!")
		log.Println("Server is listening on port 8000...")
	})

	r.Post("/media", mediaHandler.Upload)

	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}