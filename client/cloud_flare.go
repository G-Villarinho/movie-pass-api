package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/GSVillas/movie-pass-api/config"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/do"
)

type CloudFlareService interface {
	UploadImage(imageBytes []byte, filename string) (*UploadImageResponse, error)
	DeleteImage(cloudFlareID uuid.UUID) error
}

type cloudFlareService struct {
	i *do.Injector
}

type cloudflareError struct {
	Message string `json:"message"`
}

type cloudflareResult struct {
	Variants          []string `json:"variants"`
	ID                string   `json:"id"`
	Filename          string   `json:"filename"`
	Uploaded          string   `json:"uploaded"`
	RequireSignedURLs bool     `json:"requireSignedURLs"`
}

type cloudflareResponse struct {
	Messages []string          `json:"messages"`
	Success  bool              `json:"success"`
	Result   cloudflareResult  `json:"result"`
	Errors   []cloudflareError `json:"errors"`
}

type UploadImageResponse struct {
	ID  uuid.UUID
	URL string
}

func NewCloudFlareService(i *do.Injector) (CloudFlareService, error) {
	return &cloudFlareService{
		i: i,
	}, nil
}

func (c *cloudFlareService) UploadImage(imageBytes []byte, filename string) (*UploadImageResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(imageBytes)); err != nil {
		return nil, fmt.Errorf("error copying file to buffer: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	req, err := http.NewRequest("POST", config.Env.CloudFlareAccountAPI, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Env.CloudFlareApiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading API response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload error with status code: %d", resp.StatusCode)
	}

	var cloudflareResp cloudflareResponse
	if err := jsoniter.Unmarshal(respBody, &cloudflareResp); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	if !cloudflareResp.Success {
		return nil, fmt.Errorf("cloudflare response error: %+v", cloudflareResp.Errors)
	}

	imageURL := cloudflareResp.Result.Variants[0]
	imageID := cloudflareResp.Result.ID

	ID, err := uuid.Parse(imageID)
	if err != nil {
		return nil, fmt.Errorf("error to convert imageId to uuid struct: %w", err)
	}

	return &UploadImageResponse{
		ID:  ID,
		URL: imageURL,
	}, nil
}

func (c *cloudFlareService) DeleteImage(cloudFlareID uuid.UUID) error {
	deleteURL := fmt.Sprintf("%s/%s", config.Env.CloudFlareAccountAPI, cloudFlareID.String())

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Env.CloudFlareApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting image with status code: %d", resp.StatusCode)
	}

	return nil
}
