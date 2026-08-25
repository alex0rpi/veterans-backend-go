package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

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

	// media - images - videos - audio
	r.Post("/media", mediaHandler.Upload)
	/* endpoint to get a list of media items filtered by context and optionally by season query params */
	r.Get("/media", mediaHandler.ListMedia)

	/* endpoint to delete a single media element (along with its variants) */
	// r.Delete("/media/{id}", mediaHandler.Delete)

	/* endpoint to update a single media element */
	// r.Put("/media/{id}", mediaHandler.Update)

	// documents
	r.Post("/documents", documentHandler.Upload)
	r.Get("/documents", documentHandler.List)

	log.Printf("Server is listening on port %s ...", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		log.Fatal(err)
	}
}
