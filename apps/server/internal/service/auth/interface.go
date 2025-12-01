package auth

import (
	"github.com/zarazaex69/zik/apps/server/internal/config"
	"github.com/zarazaex69/zik/apps/server/internal/domain"
)

// AuthServicer defines the interface for authentication services.
type AuthServicer interface {
	GetUser(cfg *config.Config) (*domain.User, error)
	ClearCache()
}
