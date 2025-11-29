package ai

import (
	"regexp"
	"strings"

	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
	"github.com/zarazaex/zik/apps/server/internal/pkg/utils"
)

// FormatStreamChunk formats Z.AI SSE chunk to OpenAI format
func FormatStreamChunk(data map[string]interface{}, cfg *config.Config) *domain.ResponseMessage {
	// Extract nested data object
	dataObj, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil
	}

	// Get content fields
	deltaContent := getStringField(dataObj, "delta_content")
	editContent := getStringField(dataObj, "edit_content")
	phase := getStringField(dataObj, "phase")

	content := deltaContent
	if content == "" {
		content = editContent
	}
	if content == "" {
		return nil
	}

	// Apply think mode transformations based on phase
	if phase == "thinking" && cfg.Model.ThinkMode == "reasoning" {
		// Strip thinking tags for reasoning mode
		content = stripThinkingTags(content)
		return &domain.ResponseMessage{
			Role:             "assistant",
			ReasoningContent: content,
		}
	}

	// Regular content
	if content != "" {
		content = stripThinkingTags(content)
		return &domain.ResponseMessage{
			Role:    "assistant",
			Content: content,
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
func CountTokens(messages []domain.Message) int {
	text := ExtractTextFromMessages(messages)
	return utils.CountTokens(text)
}

// stripThinkingTags removes Z.AI specific thinking tags
func stripThinkingTags(content string) string {
	// Remove details tags
	content = regexp.MustCompile(`(?s)<details[^>]*?>.*?</details>`).ReplaceAllString(content, "")
	// Remove thinking tags
	content = strings.ReplaceAll(content, "<thinking>", "")
	content = strings.ReplaceAll(content, "</thinking>", "")
	// Remove summary tags
	content = regexp.MustCompile(`\n*<summary>.*?</summary>\n*`).ReplaceAllString(content, "")
	// Remove reasoning tags for reasoning mode
	content = strings.ReplaceAll(content, "<reasoning>", "")
	content = strings.ReplaceAll(content, "</reasoning>", "")
	return content
}

func getStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
