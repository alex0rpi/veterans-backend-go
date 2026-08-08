package utils

import (
	"path/filepath"

	"github.com/google/uuid"
)

func GenerateObjectKeyy(originalFilename string) string {
	extension := filepath.Ext(originalFilename)

	return "images/" + uuid.NewString() + extension
}
