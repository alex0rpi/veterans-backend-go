package config

import (
	"fmt"
	"log"
	"os"
	"veterans-go-chi-server/internal/storage"
)

func NewMediaStorage() storage.MediaStorage {
	var mediaStorage storage.MediaStorage
	STORAGE_DRIVER := os.Getenv("STORAGE_DRIVER")
		switch STORAGE_DRIVER {
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

	if(STORAGE_DRIVER == "r2") {
		fmt.Println("Connected to R2 Storage")
	} else {
		fmt.Println("Connected to Local Storage")
	}

	return mediaStorage
}

