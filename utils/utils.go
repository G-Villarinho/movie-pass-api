package utils

import (
	"fmt"
	"io"
	"mime/multipart"
)

func ConvertImageToBytes(image *multipart.FileHeader) ([]byte, error) {
	file, err := image.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %v", err)
	}

	return imageBytes, nil
}
