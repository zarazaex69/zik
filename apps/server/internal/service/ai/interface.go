package ai

import (
	"net/http"

	"github.com/zarazaex69/zik/apps/server/internal/domain"
)

// AIClienter defines the interface for the AI client.
type AIClienter interface {
	SendChatRequest(req *domain.ChatRequest, chatID string) (*http.Response, error)
}
