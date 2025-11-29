package ai

import (
	"bufio"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
	"github.com/zarazaex/zik/apps/server/internal/pkg/logger"
	"github.com/zarazaex/zik/apps/server/internal/pkg/utils"
)

var (
	phaseBak = "thinking"
)

// ParseSSEStream parses Server-Sent Events stream from Z.AI API
func ParseSSEStream(resp *http.Response) <-chan *domain.ZaiResponse {
	ch := make(chan *domain.ZaiResponse)

	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Bytes()

			// Skip empty lines or non-data lines
			if len(line) == 0 || !strings.HasPrefix(string(line), "data: ") {
				continue
			}

			// Extract JSON data
			jsonData := line[6:] // Skip "data: " prefix

			var zaiResp domain.ZaiResponse
			if err := json.Unmarshal(jsonData, &zaiResp); err != nil {
				logger.Debug().Err(err).Msg("Failed to parse SSE event")
				continue
			}

			ch <- &zaiResp
		}

		if err := scanner.Err(); err != nil {
			logger.Error().Err(err).Msg("Error reading SSE stream")
		}
	}()

	return ch
}

// FormatResponse formats Z.AI response to OpenAI format
func FormatResponse(data *domain.ZaiResponse, cfg *config.Config) map[string]interface{} {
	if data == nil || data.Data == nil {
		return nil
	}

	phase := data.Data.Phase
	if phase == "" {
		phase = "other"
	}

	content := data.Data.DeltaContent
	if content == "" {
		content = data.Data.EditContent
	}
	if content == "" {
		return nil
	}

	// Handle tool_call phase
	if phase == "tool_call" {
		content = regexp.MustCompile(`\n*<glm_block[^>]*>\{"type": "mcp", "data": \{"metadata": \{`).ReplaceAllString(content, "{")
		content = regexp.MustCompile(`[}"], "result": "".*</glm_block>`).ReplaceAllString(content, "")
	} else if phase == "other" && phaseBak == "tool_call" && strings.Contains(content, "glm_block") {
		phase = "tool_call"
		content = regexp.MustCompile(`null, "display_result": "".*</glm_block>`).ReplaceAllString(content, "\"}")
	}

	// Get config for thinking mode
	thinkMode := cfg.Model.ThinkMode

	// Handle thinking/answer phase
	if phase == "thinking" || (phase == "answer" && strings.Contains(content, "summary>")) {
		content = strings.ReplaceAll(content, "</thinking>", "")
		content = strings.ReplaceAll(content, "<Full>", "")
		content = strings.ReplaceAll(content, "</Full>", "")

		if phase == "thinking" {
			content = regexp.MustCompile(`\n*<summary>.*?</summary>\n*`).ReplaceAllString(content, "\n\n")
		}

		// Convert to <reasoning> tags
		content = regexp.MustCompile(`<details[^>]*>\n*`).ReplaceAllString(content, "<reasoning>\n\n")
		content = regexp.MustCompile(`\n*</details>`).ReplaceAllString(content, "\n\n</reasoning>")

		// Apply thinking mode transformations
		switch thinkMode {
		case "reasoning":
			if phase == "thinking" {
				content = regexp.MustCompile(`\n>\s?`).ReplaceAllString(content, "\n")
			}
			content = regexp.MustCompile(`\n*<summary>.*?</summary>\n*`).ReplaceAllString(content, "")
			content = regexp.MustCompile(`<reasoning>\n*`).ReplaceAllString(content, "")
			content = regexp.MustCompile(`\n*</reasoning>`).ReplaceAllString(content, "")

		case "think":
			if phase == "thinking" {
				content = regexp.MustCompile(`\n>\s?`).ReplaceAllString(content, "\n")
			}
			content = regexp.MustCompile(`\n*<summary>.*?</summary>\n*`).ReplaceAllString(content, "")
			content = strings.ReplaceAll(content, "<reasoning>", "<think>")
			content = strings.ReplaceAll(content, "</reasoning>", "</think>")

		case "strip":
			content = regexp.MustCompile(`\n*<summary>.*?</summary>\n*`).ReplaceAllString(content, "")
			content = regexp.MustCompile(`<reasoning>\n*`).ReplaceAllString(content, "")
			content = regexp.MustCompile(`</reasoning>`).ReplaceAllString(content, "")

		default:
			content = regexp.MustCompile(`</reasoning>`).ReplaceAllString(content, "</reasoning>\n\n")
		}
	}

	phaseBak = phase

	// Return formatted response based on type
	if phase == "thinking" && thinkMode == "reasoning" {
		return map[string]interface{}{
			"role":              "assistant",
			"reasoning_content": content,
		}
	}

	if phase == "tool_call" {
		// Note: The original proxy returns "tool_call": content
		// We keep it as map to let handler decide how to process it
		// Or we can try to parse it if it's a complete JSON, but for streaming it might be partial.
		// Z.AI sends tool calls as text chunks that accumulate to JSON.
		return map[string]interface{}{
			"tool_call": content,
		}
	}

	if content != "" {
		return map[string]interface{}{
			"role":    "assistant",
			"content": content,
		}
	}

	return nil
}

// ExtractTextFromMessages extracts text content from messages for token counting
func ExtractTextFromMessages(messages []domain.Message) string {
	var texts []string

	for _, msg := range messages {
		// Handle string content
		if contentStr, ok := msg.Content.(string); ok {
			texts = append(texts, contentStr)
			continue
		}

		// Handle array content (multimodal)
		if contentArr, ok := msg.Content.([]interface{}); ok {
			for _, item := range contentArr {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if itemMap["type"] == "text" {
						if text, ok := itemMap["text"].(string); ok {
							texts = append(texts, text)
						}
					}
				}
			}
		}
	}

	return strings.Join(texts, " ")
}

// CountTokens counts tokens in messages
func CountTokens(messages []domain.Message, tokenizer utils.Tokener) int {
	text := ExtractTextFromMessages(messages)
	return tokenizer.Count(text)
}
