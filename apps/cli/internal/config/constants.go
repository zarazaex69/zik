package config

import "os"

// API endpoint is hardcoded at compile time to prevent users from redirecting requests
// to unauthorized servers. This ensures all requests go through the official ZIK API.
const (
	// DefaultAPIEndpoint is the production API endpoint that cannot be overridden by users
	DefaultAPIEndpoint = "https://api.zik.zarazaex.xyz"

	// DefaultModel is the default AI model to use for requests
	DefaultModel = "GLM-4-6-API-V1"

	// Version is the CLI version, can be overridden at build time with -ldflags
	Version = "0.1.0"
)

// GetAPIEndpoint returns the API endpoint, allowing override via ZIK_API_URL env var for development
// In production builds, this env var check can be removed for security
func GetAPIEndpoint() string {
	// Allow override for development/testing purposes
	if endpoint := os.Getenv("ZIK_API_URL"); endpoint != "" {
		return endpoint
	}
	return DefaultAPIEndpoint
}
