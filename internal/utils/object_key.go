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

	return "images/" + id + "/" + string(constants.VariantOriginal) + extension
}

func GenerateVariantObjectKey(objectKey, variantName string) string {
	directory := path.Dir(objectKey)

	return path.Join(directory, variantName+".webp")
}

// Documents share the same bucket as images but live under their own prefix, since they have no variants.
func GenerateDocumentObjectKey(originalFilename string) string {
	id := uuid.NewString()

	extension := strings.ToLower(filepath.Ext(originalFilename))

	return "documents/" + id + extension
}
