package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
	"github.com/zarazaex/zik/apps/server/internal/pkg/logger"
	"github.com/zarazaex/zik/apps/server/internal/pkg/utils"
	"github.com/zarazaex/zik/apps/server/internal/service/auth"
)

// FormatRequest converts OpenAI format to Z.AI format
func FormatRequest(req *domain.ChatRequest, cfg *config.Config) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Get model
	model := req.Model
	if model == "" {
		model = cfg.Model.Default
	}

	// Process messages
	newMessages := []map[string]interface{}{}
	chatID := utils.GenerateRequestID()

	for _, msg := range req.Messages {
		newMessage := map[string]interface{}{"role": msg.Role}

		// Handle string content
		if contentStr, ok := msg.Content.(string); ok {
			newMessage["content"] = contentStr
			newMessages = append(newMessages, newMessage)
			continue
		}

		// Handle array content (multimodal)
		if contentArr, ok := msg.Content.([]interface{}); ok {
			var newContent interface{} = ""

			for _, item := range contentArr {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				itemType, _ := itemMap["type"].(string)

				// Text content
				if itemType == "text" {
					if text, ok := itemMap["text"].(string); ok {
						newContent = text
					}
					continue
				}

				// Image content
				if itemType == "image_url" {
					mediaURL := ""

					// OpenAI format
					if imageURL, ok := itemMap["image_url"].(map[string]interface{}); ok {
						if urlStr, ok := imageURL["url"].(string); ok {
							mediaURL = urlStr
						}
					}

					if mediaURL == "" {
						continue
					}

					// Upload image if it's base64
					uploadedURL, err := UploadImage(mediaURL, chatID, cfg)
					if err != nil {
						logger.Warn().Err(err).Msg("Failed to upload image")
						continue
					}
					if uploadedURL != "" {
						mediaURL = uploadedURL
					}

					// Convert newContent to array if needed
					if contentStr, ok := newContent.(string); ok {
						newContent = []map[string]interface{}{
							{"type": "text", "text": contentStr},
						}
					}
					if contentSlice, ok := newContent.([]map[string]interface{}); ok {
						contentSlice = append(contentSlice, map[string]interface{}{
							"type":      "image_url",
							"image_url": map[string]interface{}{"url": mediaURL},
						})
						newContent = contentSlice
					}
				}
			}

			newMessage["content"] = newContent
			newMessages = append(newMessages, newMessage)
		}
	}

	result["model"] = model
	result["messages"] = newMessages
	result["stream"] = true

	// Handle features
	features := map[string]interface{}{
		"enable_thinking": false,
	}

	if len(features) > 0 {
		result["features"] = features
	}

	return result, nil
}

// UploadImage uploads a base64 image to Z.AI API
func UploadImage(dataURL, chatID string, cfg *config.Config) (string, error) {
	// Skip upload in anonymous mode or if not base64
	if cfg.Upstream.Anonymous || !strings.HasPrefix(dataURL, "data:") {
		return "", nil
	}

	// Parse data URL
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid data URL format")
	}

	encoded := parts[1]

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Generate filename
	filename := utils.GenerateID()

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(imageData)); err != nil {
		return "", fmt.Errorf("failed to write file data: %w", err)
	}
	writer.Close()

	// Get user token
	authService := auth.NewService()
	user, err := authService.GetUser(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	// Build request
	uploadURL := fmt.Sprintf("%s//%s/api/v1/files/", cfg.Upstream.Protocol, cfg.Upstream.Host)
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	// Set headers
	headers := cfg.GetUpstreamHeaders()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.Token))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Referer", fmt.Sprintf("%s//%s/c/%s", cfg.Upstream.Protocol, cfg.Upstream.Host, chatID))

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var result struct {
		ID       string `json:"id"`
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse upload response: %w", err)
	}

	return fmt.Sprintf("%s_%s", result.ID, result.Filename), nil
}
