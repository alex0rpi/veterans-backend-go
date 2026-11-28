package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/crypto/bcrypt"

	"veterans-go-chi-server/internal/auth"
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be set and at least 32 characters long")
	}
	jwtManager := auth.NewJWTManager(jwtSecret)

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPasswordHash := os.Getenv("ADMIN_PASSWORD_HASH")

	if adminUsername == "" || adminPasswordHash == "" {
		log.Fatal("ADMIN_USERNAME and ADMIN_PASSWORD_HASH must be set")
	}

	if _, err := bcrypt.Cost([]byte(adminPasswordHash)); err != nil {
    log.Fatal("ADMIN_PASSWORD_HASH must be a valid bcrypt hash")
}

	authHandler := handlers.NewAuthHandler(jwtManager, adminUsername, adminPasswordHash)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// IMAGES ENDPOINTS //

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintln(w, "Bon dia desde Go Chi!"); err != nil {
			log.Println("Error writing response:", err)
		}
	})

	//* Media public endpoints -------

	/* endpoint to get a list of media items filtered by context and optionally by season and category query params */
	r.Get("/media", mediaHandler.ListMedia)

	// AUTH ENDPOINTS //

	r.Post("/auth/login", authHandler.Login)
	// r.Post("/auth/logout", authHandler.Logout)

	//* Admin protected endpoints -------

	r.Group(func(r chi.Router) {
		r.Use(jwtManager.AuthMiddleware)

		// IMAGE ENDPOINTS //
		/* endpoint to post a new media item (image) */
		r.Post("/media", mediaHandler.UploadImage)
		/* endpoint to delete a single media element (along with its variants) */
		r.Delete("/media/{id}", mediaHandler.DeleteImageById)
		/* endpoint to get a single image element metadata by its ID */
		r.Get("/media/{id}", mediaHandler.GetImageMetadataById)
		/* endpoint to update a single media element */
		r.Patch("/media/{id}", mediaHandler.UpdateImage)
		/* endpoint to update a images display position for a given context, season and category */
		r.Patch("/media/display-positions", mediaHandler.UpdateImageDisplayPositions)

		// DOCUMENTS ENDPOINTS //
		/* endpoint to post a new document */
		r.Post("/documents", documentHandler.Upload)
		/* endpoint to get a list of all documents metadata */
		r.Get("/documents", documentHandler.List)

	})

	// Start the HTTP server

	log.Printf("Server is UP and listening on port %s ...", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		log.Fatal(err)
	}
}
