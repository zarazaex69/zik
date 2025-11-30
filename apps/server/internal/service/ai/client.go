package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
	"github.com/zarazaex/zik/apps/server/internal/pkg/crypto"
	"github.com/zarazaex/zik/apps/server/internal/pkg/httpclient"
	"github.com/zarazaex/zik/apps/server/internal/pkg/logger"
	"github.com/zarazaex/zik/apps/server/internal/pkg/utils"
	"github.com/zarazaex/zik/apps/server/internal/service/auth"
)

// Client handles communication with Z.AI API
type Client struct {
	cfg          *config.Config
	authService  auth.AuthServicer
	signatureGen crypto.SignatureGenerator // New field for dependency injection
}

// NewClient creates a new Z.AI API client
func NewClient(cfg *config.Config, authSvc auth.AuthServicer, sigGen crypto.SignatureGenerator) *Client {
	return &Client{
		cfg:          cfg,
		authService:  authSvc,
		signatureGen: sigGen, // Assign the injected signature generator
	}
}

// SendChatRequest sends a chat completion request to Z.AI API
func (c *Client) SendChatRequest(req *domain.ChatRequest, chatID string) (*http.Response, error) {
	timestamp := time.Now().UnixMilli()
	requestID := utils.GenerateRequestID()

	// Get user info for authentication
	user, err := c.authService.GetUser(c.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Build query parameters
	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", timestamp))
	params.Set("requestId", requestID)
	params.Set("version", "0.0.1")
	params.Set("platform", "web")
	params.Set("token", user.Token)

	// Build headers
	headers := c.cfg.GetUpstreamHeaders()
	headers["Authorization"] = "Bearer " + user.Token
	headers["Content-Type"] = "application/json"
	headers["Referer"] = fmt.Sprintf("%s//%s/c/%s", c.cfg.Upstream.Protocol, c.cfg.Upstream.Host, chatID)

	// Prepare request body using FormatRequest
	requestBody, err := FormatRequest(req, c.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to format request: %w", err)
	}

	// Add required fields
	requestBody["chat_id"] = chatID
	requestBody["id"] = utils.GenerateRequestID()

	// Add signature for authenticated users
	if user.ID != "" {
		params.Set("user_id", user.ID)

		// Extract last user message for signature
		lastUserMsg := extractLastUserMessage(req.Messages)

		// Generate signature using the injected dependency
		sigParams := map[string]string{
			"requestId": requestID,
			"timestamp": fmt.Sprintf("%d", timestamp),
			"user_id":   user.ID,
		}
		sigResult, err := c.signatureGen.GenerateSignature(sigParams, lastUserMsg) // Use injected interface
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to generate signature, continuing without it")
		} else {
			headers["x-signature"] = sigResult.Signature
			params.Set("signature_timestamp", fmt.Sprintf("%d", sigResult.Timestamp))
			requestBody["signature_prompt"] = lastUserMsg
		}
	}

	// Build URL
	apiURL := fmt.Sprintf("%s//%s/api/v2/chat/completions?%s",
		c.cfg.Upstream.Protocol, c.cfg.Upstream.Host, params.Encode())

	// Marshal request body
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	logger.Debug().
		Str("url", apiURL).
		Str("chat_id", chatID).
		Str("model", req.Model).
		Str("body", string(bodyBytes)).
		Msg("Sending chat request to Z.AI")

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	// Send request with no timeout for streaming
	client := httpclient.New(0)
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, domain.NewUpstreamError(resp.StatusCode, "Z.AI API error")
	}

	return resp, nil
}

// extractLastUserMessage extracts the last user message content for signature
func extractLastUserMessage(messages []domain.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			// Handle string content
			if contentStr, ok := messages[i].Content.(string); ok {
				return contentStr
			}

			// Handle array content (multimodal)
			if contentArr, ok := messages[i].Content.([]interface{}); ok {
				var texts []string
				for _, item := range contentArr {
					if itemMap, ok := item.(map[string]interface{}); ok {
						if itemMap["type"] == "text" {
							if text, ok := itemMap["text"].(string); ok {
								texts = append(texts, text)
							}
						}
					}
				}
				return strings.Join(texts, " ")
			}
		}
	}
	return ""
}
