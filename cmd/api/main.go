package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	config "veterans-go-chi-server/internal/config"
	"veterans-go-chi-server/internal/handlers"
	"veterans-go-chi-server/internal/repositories"
	"veterans-go-chi-server/internal/services"
)

func main() {
	pool, err := config.NewPostgresPool()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()


	mediaStorage := config.NewMediaStorage()

	mediaRepository := repositories.NewMediaRepository(pool)
	mediaService := services.NewMediaService(mediaRepository, mediaStorage, os.Getenv("R2_PUBLIC_BASE_URL"))
	mediaHandler := handlers.NewMediaHandler(mediaService)

	documentRepository := repositories.NewDocumentRepository(pool)
	documentService := services.NewDocumentService(documentRepository, mediaStorage, os.Getenv("R2_PUBLIC_BASE_URL"))
	documentHandler := handlers.NewDocumentHandler(documentService)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintln(w, "Bon dia desde Go Chi!"); err != nil {
			log.Println("Error writing response:", err)
		}
	})

	/* endpoint to post a new media item (image) */
	r.Post("/media", mediaHandler.Upload)

	/* endpoint to get a list of media items filtered by context and optionally by season and category query params */
	r.Get("/media", mediaHandler.ListMedia)

	// TODO *********************************************************************************
	/* endpoint to get a single media element metadata by its ID */
	// r.Get("/media/{id}", mediaHandler.GetMediaById)

	/* endpoint to update a single media element */
	// r.Put("/media/{id}", mediaHandler.Update)

	/* endpoint to delete a single media element (along with its variants) */
	r.Delete("/media/{id}", mediaHandler.DeleteImage)

	/* endpoint to get all available display positions for a given context and category */
	r.Get("/media/available-display-positions", mediaHandler.GetAvailableDisplayPositions)

	// documents
	/* endpoint to post a new document */
	r.Post("/documents", documentHandler.Upload)

	/* endpoint to get a list of all documents metadata */
	r.Get("/documents", documentHandler.List)

	log.Printf("Server is UP and listening on port %s ...", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		log.Fatal(err)
	}
}
