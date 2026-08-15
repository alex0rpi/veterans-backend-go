package imaging

import (
	"bytes"
	"image"
	"veterans-go-chi-server/internal/constants"

	"github.com/deepteams/webp"
	"github.com/disintegration/imaging"
)

type Variant struct {
	Name			string
	MaxDimension	int
}

// Slice of image variants that can be used for resizing images.
var Variants = []Variant{
	{ Name: constants.VariantBlur, MaxDimension: 32},
	{ Name: constants.VariantSmall, MaxDimension: 480},
	{ Name: constants.VariantMedium, MaxDimension: 1024},
	{ Name: constants.VariantLarge, MaxDimension: 1920},
}

type GeneratedVariant struct {
	Name			string
	Data			[]byte
}

func GenerateVariants(img image.Image) ([]GeneratedVariant, error) {
	var generatedVariants []GeneratedVariant
	for _, variant := range Variants {

		//* Calculate dimensions
		newWidth, newHeight := calculateVariantSize(
			img.Bounds().Dx(),
			 img.Bounds().Dy(),
			  variant.MaxDimension,
			)

		//* Do perform the actual resizing and encoding of the image
		resizedImg := resizeImage(img, newWidth, newHeight)

		/* if variant.Name == "blur" {
			resizedImg = imaging.Blur(resizedImg, 2)
		} */

		imgData, err := encodeToWebP(resizedImg, newWidth, newHeight)
		if err != nil {
			return nil, err
		}

		generatedVariants = append(generatedVariants, GeneratedVariant{
			Name: variant.Name,
			Data: imgData,
		})
	}
	return generatedVariants, nil

}

func encodeToWebP(img image.Image, width, height int) ([]byte, error) {
	var buffer bytes.Buffer

	err := webp.Encode(&buffer, img, &webp.EncoderOptions{
		Quality: 80,
		Method: 4,
	})
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func calculateVariantSize( width int, height int, maxDimension int) (int, int) {
	if width <= maxDimension && height <= maxDimension {
		return width, height
	}
	// If not, calculate the new dimensions while maintaining the aspect ratio
	if width > height {
		newWidth := maxDimension // set the biggest dimension to maxDimension
		newHeight := height * maxDimension / width // proportional calculation to adapt the other dimension
		return newWidth, newHeight
	} else {
		newHeight := maxDimension
		newWidth := width * maxDimension / height
		return newWidth, newHeight
	}
}

func resizeImage(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}