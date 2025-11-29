package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// SignatureResult holds the generated signature and timestamp.
type SignatureResult struct {
	Signature string
	Timestamp int64
}

// SignatureGenerator defines the interface for generating signatures.
type SignatureGenerator interface {
	GenerateSignature(params map[string]string, lastUserMessage string) (*SignatureResult, error)
}

// defaultSignatureGenerator is the default implementation of SignatureGenerator.
type defaultSignatureGenerator struct{}

// NewSignatureGenerator creates a new default SignatureGenerator.
func NewSignatureGenerator() SignatureGenerator {
	return &defaultSignatureGenerator{}
}

// GenerateSignature generates a signature for the given parameters and last user message.
func (s *defaultSignatureGenerator) GenerateSignature(params map[string]string, lastUserMessage string) (*SignatureResult, error) {
	// Construct the string to be signed
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
		sb.WriteString("&")
	}

	sb.WriteString(lastUserMessage)
	stringToSign := sb.String()

	// Get secret key from environment variable
	// For production, use a more secure method to retrieve the secret key
	secretKey := os.Getenv("ZAI_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("ZAI_SECRET_KEY environment variable not set")
	}

	// Generate HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h.Sum(nil))

	return &SignatureResult{
		Signature: signature,
		Timestamp: time.Now().UnixMilli(),
	}, nil
}

