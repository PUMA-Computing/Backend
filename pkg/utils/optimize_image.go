package utils

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"io"
)

func OptimizeImage(src io.Reader, maxWidth, maxHeight int) (io.Reader, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	resizedImg := imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)

	// Create a buffer to store the encoded image
	buf := new(bytes.Buffer)

	// Encode the image to WEBP format
	err = jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, err
	}

	// Convert the buffer to io.Reader and return
	return bytes.NewReader(buf.Bytes()), nil
}
