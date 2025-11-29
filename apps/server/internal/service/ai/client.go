package ai

import (
	"bufio"
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
	"github.com/zarazaex/zik/apps/server/internal/pkg/logger"
	"github.com/zarazaex/zik/apps/server/internal/pkg/utils"
	"github.com/zarazaex/zik/apps/server/internal/service/auth"
)

// Client handles communication with Z.AI API
type Client struct {
	cfg         *config.Config
	authService *auth.Service
}

// NewClient creates a new Z.AI API client
func NewClient(cfg *config.Config, authSvc *auth.Service) *Client {
	return &Client{
		cfg:         cfg,
		authService: authSvc,
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

	// Build headers
	headers := c.cfg.GetUpstreamHeaders()
	headers["Authorization"] = "Bearer " + user.Token
	headers["Content-Type"] = "application/json"
	headers["Referer"] = fmt.Sprintf("%s//%s/c/%s", c.cfg.Upstream.Protocol, c.cfg.Upstream.Host, chatID)

	// Prepare request body
	requestBody := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   true,
		"chat_id":  chatID,
		"id":       utils.GenerateRequestID(),
	}

	// Add signature for authenticated users
	if user.ID != "" {
		params.Set("user_id", user.ID)

		// Extract last user message for signature
		lastUserMsg := extractLastUserMessage(req.Messages)

		// Generate signature
		sigParams := map[string]string{
			"requestId": requestID,
			"timestamp": fmt.Sprintf("%d", timestamp),
			"user_id":   user.ID,
		}
		sigResult, err := crypto.GenerateSignature(sigParams, lastUserMsg)
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to generate signature, continuing without it")
		} else {
			headers["X-Signature"] = sigResult.Signature
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
	client := &http.Client{Timeout: 0}
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

// ParseSSEStream parses Server-Sent Events stream from Z.AI API
func ParseSSEStream(resp *http.Response) <-chan map[string]interface{} {
	ch := make(chan map[string]interface{})

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// Skip empty lines or non-data lines
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			// Extract JSON data
			jsonData := line[6:] // Skip "data: " prefix

			var data map[string]interface{}
			if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
				logger.Debug().Err(err).Msg("Failed to parse SSE event")
				continue
			}

			ch <- data
		}

		if err := scanner.Err(); err != nil {
			logger.Error().Err(err).Msg("Error reading SSE stream")
		}
	}()

	return ch
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
