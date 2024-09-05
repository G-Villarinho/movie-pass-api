package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/GSVillas/movie-pass-api/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/do"
)

var (
	ErrCreateFormFile   = errors.New("failed to create form file")
	ErrCopyFile         = errors.New("failed to copy file to buffer")
	ErrCloseWriter      = errors.New("failed to close writer")
	ErrCreateRequest    = errors.New("failed to create request")
	ErrSendRequest      = errors.New("failed to send request")
	ErrReadResponse     = errors.New("failed to read API response")
	ErrDecodeJSON       = errors.New("failed to decode JSON response")
	ErrUploadFailed     = errors.New("upload failed with status code")
	ErrCloudflareFailed = errors.New("cloudflare response error")
)

type CloudFlareService interface {
	UploadImage(imageBytes []byte, filename string) (string, error)
}

type cloudFlareService struct {
	i *do.Injector
}

type CloudflareError struct {
	Message string `json:"message"`
}

type CloudflareResult struct {
	Variants          []string `json:"variants"`
	ID                string   `json:"id"`
	Filename          string   `json:"filename"`
	Uploaded          string   `json:"uploaded"`
	RequireSignedURLs bool     `json:"requireSignedURLs"`
}

type CloudflareResponse struct {
	Messages []string          `json:"messages"`
	Success  bool              `json:"success"`
	Result   CloudflareResult  `json:"result"`
	Errors   []CloudflareError `json:"errors"`
}

func NewCloudFlareService(i *do.Injector) (CloudFlareService, error) {
	return &cloudFlareService{
		i: i,
	}, nil
}

func (c *cloudFlareService) UploadImage(imageBytes []byte, filename string) (string, error) {
	log := slog.With(
		slog.String("client", "cloudFlare"),
		slog.String("func", "UploadImage"),
	)

	log.Info("Initializing image upload process")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Cria o arquivo de formulário
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Error("Failed to create form file", slog.String("error", err.Error()))
		return "", ErrCreateFormFile
	}

	// Copia o array de bytes para o buffer
	if _, err := io.Copy(part, bytes.NewReader(imageBytes)); err != nil {
		log.Error("Failed to copy file to buffer", slog.String("error", err.Error()))
		return "", ErrCopyFile
	}

	// Fecha o writer
	if err := writer.Close(); err != nil {
		log.Error("Failed to close writer", slog.String("error", err.Error()))
		return "", ErrCloseWriter
	}

	req, err := http.NewRequest("POST", config.Env.CloudFlareAccountAPI, body)
	if err != nil {
		log.Error("Failed to create request", slog.String("error", err.Error()))
		return "", ErrCreateRequest
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Env.CloudFlareApiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to send request", slog.String("error", err.Error()))
		return "", ErrSendRequest
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read API response", slog.String("error", err.Error()))
		return "", ErrReadResponse
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Upload failed", slog.Int("status", resp.StatusCode))
		return "", ErrUploadFailed
	}

	var cloudflareResp CloudflareResponse
	if err := jsoniter.Unmarshal(respBody, &cloudflareResp); err != nil {
		log.Error("Failed to decode JSON response", slog.String("error", err.Error()))
		return "", ErrDecodeJSON
	}

	if !cloudflareResp.Success {
		log.Error("Cloudflare response error", slog.String("error", ErrCloudflareFailed.Error()), slog.Any("details", cloudflareResp.Errors))
		return "", ErrCloudflareFailed
	}

	imageURL := cloudflareResp.Result.Variants[0]

	log.Info("Image upload successful", slog.String("imageURL", imageURL))
	return imageURL, nil
}
