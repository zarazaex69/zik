package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zarazaex69/zik/apps/server/internal/config"
	"github.com/zarazaex69/zik/apps/server/internal/domain"
	"github.com/zarazaex69/zik/apps/server/internal/pkg/logger"
	"github.com/zarazaex69/zik/apps/server/internal/pkg/utils"
	"github.com/zarazaex69/zik/apps/server/internal/pkg/validator"
	"github.com/zarazaex69/zik/apps/server/internal/service/ai"
)

// ChatCompletions handles OpenAI-compatible chat completions
func ChatCompletions(cfg *config.Config, aiClient ai.AIClienter, tokenizer utils.Tokener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate request
		var req domain.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON request")
			return
		}

		// Validate request
		if err := validator.Validate(&req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Set default model if not specified
		if req.Model == "" {
			req.Model = cfg.Model.Default
		}

		// Generate chat ID
		chatID := utils.GenerateRequestID()

		logger.Info().
			Str("model", req.Model).
			Bool("stream", req.Stream).
			Int("messages", len(req.Messages)).
			Msg("Processing chat completion request")

		// Send request to Z.AI
		resp, err := aiClient.SendChatRequest(&req, chatID)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to send chat request")
			writeError(w, http.StatusInternalServerError, "Failed to process request")
			return
		}

		// Handle streaming response
		if req.Stream {
			handleStreamingResponse(w, resp, &req, cfg, tokenizer)
		} else {
			handleNonStreamingResponse(w, resp, &req, cfg, tokenizer)
		}
	}
}

func handleStreamingResponse(w http.ResponseWriter, resp *http.Response, req *domain.ChatRequest, cfg *config.Config, tokenizer utils.Tokener) {
	// Check flusher support BEFORE setting headers
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "Streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Track content for usage calculation
	var contentParts []string
	includeUsage := req.StreamOpts != nil && req.StreamOpts.IncludeUsage

	// Calculate prompt tokens
	promptTokens := 0
	if includeUsage {
		promptTokens = ai.CountTokens(req.Messages, tokenizer)
	}

	// Stream responses
	for zaiResp := range ai.ParseSSEStream(resp) {
		delta := ai.FormatResponse(zaiResp, cfg)
		if delta == nil {
			continue
		}

		// Collect content for token counting
		if includeUsage {
			if content, ok := delta["content"].(string); ok {
				contentParts = append(contentParts, content)
			}
			if reasoningContent, ok := delta["reasoning_content"].(string); ok {
				contentParts = append(contentParts, reasoningContent)
			}
		}

		// Send chunk
		deltaResponse := &domain.ResponseMessage{
			Role:             getStringFromMap(delta, "role"),
			Content:          getStringFromMap(delta, "content"),
			ReasoningContent: getStringFromMap(delta, "reasoning_content"),
			ToolCall:         getStringFromMap(delta, "tool_call"),
		}

		streamChunk := domain.ChatResponse{
			ID:      utils.GenerateChatCompletionID(),
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []domain.Choice{
				{
					Index: 0,
					Delta: deltaResponse,
				},
			},
		}

		chunkJSON, _ := json.Marshal(streamChunk)
		fmt.Fprintf(w, "data: %s\n\n", chunkJSON)
		flusher.Flush()
	}

	// Send finish reason
	finishChunk := domain.ChatResponse{
		ID:      utils.GenerateChatCompletionID(),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []domain.Choice{
			{
				Index:        0,
				Delta:        &domain.ResponseMessage{Role: "assistant"},
				FinishReason: stringPtr("stop"),
			},
		},
	}
	finishJSON, _ := json.Marshal(finishChunk)
	fmt.Fprintf(w, "data: %s\n\n", finishJSON)
	flusher.Flush()

	// Send usage if requested
	if includeUsage {
		completionText := strings.Join(contentParts, "")
		completionTokens := tokenizer.Count(completionText)

		usageChunk := domain.ChatResponse{
			ID:      utils.GenerateChatCompletionID(),
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []domain.Choice{},
			Usage: &domain.Usage{
				PromptTokens:     promptTokens,
				CompletionTokens: completionTokens,
				TotalTokens:      promptTokens + completionTokens,
			},
		}
		usageJSON, _ := json.Marshal(usageChunk)
		fmt.Fprintf(w, "data: %s\n\n", usageJSON)
		flusher.Flush()
	}

	// Send [DONE]
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

func handleNonStreamingResponse(w http.ResponseWriter, resp *http.Response, req *domain.ChatRequest, cfg *config.Config, tokenizer utils.Tokener) {
	var contentParts []string
	var reasoningParts []string

	// Collect all chunks
	for zaiResp := range ai.ParseSSEStream(resp) {
		if zaiResp.Data != nil && zaiResp.Data.Done {
			break
		}

		delta := ai.FormatResponse(zaiResp, cfg)
		if delta == nil {
			continue
		}

		if content, ok := delta["content"].(string); ok {
			contentParts = append(contentParts, content)
		}
		if reasoningContent, ok := delta["reasoning_content"].(string); ok {
			reasoningParts = append(reasoningParts, reasoningContent)
		}
	}

	// Build final message
	finalMessage := &domain.ResponseMessage{
		Role: "assistant",
	}

	completionText := ""
	if len(reasoningParts) > 0 {
		reasoning := strings.Join(reasoningParts, "")
		finalMessage.ReasoningContent = reasoning
		completionText += reasoning
	}
	if len(contentParts) > 0 {
		content := strings.Join(contentParts, "")
		finalMessage.Content = content
		completionText += content
	}

	// Build response
	response := domain.ChatResponse{
		ID:      utils.GenerateChatCompletionID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []domain.Choice{
			{
				Index:        0,
				Message:      finalMessage,
				FinishReason: stringPtr("stop"),
			},
		},
	}

	// Add usage
	promptTokens := ai.CountTokens(req.Messages, tokenizer)
	completionTokens := tokenizer.Count(completionText)
	response.Usage = &domain.Usage{
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(domain.NewAPIError(code, message))
}

func stringPtr(s string) *string {
	return &s
}
