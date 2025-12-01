package ai

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarazaex69/zik/apps/server/internal/config"
	"github.com/zarazaex69/zik/apps/server/internal/domain"
)

func TestParseSSEStream(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string // Expected phases or content to verify parsing
	}{
		{
			name: "Valid stream",
			input: `data: {"data": {"phase": "thinking", "delta_content": "thinking..."}}

data: {"data": {"phase": "answer", "delta_content": "answer"}}
`,
			expected: []string{"thinking", "answer"},
		},
		{
			name: "Empty lines and invalid data",
			input: `
data: {"data": {"phase": "thinking"}}

invalid line
data: invalid json

data: {"data": {"phase": "answer"}}
`,
			expected: []string{"thinking", "answer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response body
			body := io.NopCloser(bytes.NewBufferString(tt.input))
			resp := &http.Response{
				Body: body,
			}

			ch := ParseSSEStream(resp)

			var phases []string
			for msg := range ch {
				if msg.Data != nil {
					phases = append(phases, msg.Data.Phase)
				}
			}

			assert.Equal(t, len(tt.expected), len(phases))
			assert.Equal(t, tt.expected, phases)
		})
	}
}

func TestFormatResponse(t *testing.T) {
	baseConfig := &config.Config{
		Model: config.ModelConfig{
			ThinkMode: "reasoning",
		},
	}

	tests := []struct {
		name      string
		input     *domain.ZaiResponse
		cfg       *config.Config
		want      map[string]interface{}
		checkFunc func(t *testing.T, got map[string]interface{})
	}{
		{
			name:  "Nil input",
			input: nil,
			cfg:   baseConfig,
			want:  nil,
		},
		{
			name: "Empty content",
			input: &domain.ZaiResponse{
				Data: &domain.ZaiResponseData{
					Phase:        "answer",
					DeltaContent: "",
				},
			},
			cfg:  baseConfig,
			want: nil,
		},
		{
			name: "Simple answer",
			input: &domain.ZaiResponse{
				Data: &domain.ZaiResponseData{
					Phase:        "answer",
					DeltaContent: "Hello",
				},
			},
			cfg: baseConfig,
			want: map[string]interface{}{
				"role":    "assistant",
				"content": "Hello",
			},
		},
		{
			name: "Thinking phase with reasoning mode",
			input: &domain.ZaiResponse{
				Data: &domain.ZaiResponseData{
					Phase:        "thinking",
					DeltaContent: "\n> thinking step",
				},
			},
			cfg: &config.Config{Model: config.ModelConfig{ThinkMode: "reasoning"}},
			want: map[string]interface{}{
				"role":              "assistant",
				"reasoning_content": "\nthinking step",
			},
		},
		{
			name: "Thinking phase with think mode",
			input: &domain.ZaiResponse{
				Data: &domain.ZaiResponseData{
					Phase:        "thinking",
					DeltaContent: "<details>thought</details>",
				},
			},
			cfg: &config.Config{Model: config.ModelConfig{ThinkMode: "think"}},
			checkFunc: func(t *testing.T, got map[string]interface{}) {
				content, ok := got["content"].(string)
				assert.True(t, ok)
				assert.Contains(t, content, "<think>")
				assert.Contains(t, content, "thought")
				assert.Contains(t, content, "</think>")
			},
		},
		{
			name: "Tool call phase",
			input: &domain.ZaiResponse{
				Data: &domain.ZaiResponseData{
					Phase:        "tool_call",
					DeltaContent: `<glm_block>{"type": "mcp", "data": {"metadata": {"name": "test"}, "result": ""}}</glm_block>`,
				},
			},
			cfg: baseConfig,
			checkFunc: func(t *testing.T, got map[string]interface{}) {
				content, ok := got["tool_call"].(string)
				assert.True(t, ok)
				assert.Contains(t, content, `"name": "test"`)
				assert.NotContains(t, content, "glm_block")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatResponse(tt.input, tt.cfg)
			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExtractTextFromMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []domain.Message
		want     string
	}{
		{
			name: "String content",
			messages: []domain.Message{
				{Content: "Hello"},
				{Content: "World"},
			},
			want: "Hello World",
		},
		{
			name: "Multimodal content",
			messages: []domain.Message{
				{
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Image"},
						map[string]interface{}{"type": "image_url", "image_url": "http://example.com/img.png"},
					},
				},
				{Content: "Description"},
			},
			want: "Image Description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTextFromMessages(tt.messages)
			assert.Equal(t, tt.want, got)
		})
	}
}
