package utils

import (
	"path"
	"path/filepath"
	"strings"
	"veterans-go-chi-server/internal/constants"

	"github.com/google/uuid"
)

func GenerateObjectKey(originalFilename string) string {

	id := uuid.NewString()

	extension := strings.ToLower(filepath.Ext(originalFilename))

	return id + "/" + constants.VariantOriginal + extension
}

func GenerateVariantObjectKey(objectKey, variantName string) string {
	directory := path.Dir(objectKey)

	return path.Join(directory, variantName + ".webp")
}
