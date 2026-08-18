package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"veterans-go-chi-server/internal/database"
	"veterans-go-chi-server/internal/handlers"
	"veterans-go-chi-server/internal/repositories"
	"veterans-go-chi-server/internal/services"
	"veterans-go-chi-server/internal/storage"
)

func main() {
	pool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	var mediaStorage storage.MediaStorage

	switch os.Getenv("STORAGE_DRIVER") {
	case "r2":
		mediaStorage = storage.NewR2Storage(
			os.Getenv("R2_ACCOUNT_ID"),
			os.Getenv("R2_ACCESS_KEY_ID"),
			os.Getenv("R2_SECRET_ACCESS_KEY"),
			os.Getenv("R2_BUCKET"),
		)
	case "local", "":
		mediaStorage = storage.NewLocalStorage("./storage")
	default:
		log.Fatalf("unsupported storage driver")
	}
	mediaRepository := repositories.NewMediaRepository(pool)
	mediaService := services.NewMediaService(mediaRepository, mediaStorage)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintln(w, "Bon dia desde Go Chi!"); err != nil {
			log.Println("Error writing response:", err)
		}
	})

	r.Post("/media", mediaHandler.Upload)
	/* endpoint to get a list of media items filtered by season and category using query params */
	// r.Get("/media", mediaHandler.List)

	/* endpoint to delete a single media element (along with its variants) */
	// r.Delete("/media/{id}", mediaHandler.Delete)

	/* endpoint to update a single media element */
	// r.Put("/media/{id}", mediaHandler.Update)

	log.Printf("Server is listening on port %s ...", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		log.Fatal(err)
	}
}
