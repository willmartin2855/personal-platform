package secret

import (
	"fmt"
	"os"
)

// SecretStore abstracts secret retrieval. Swap in an AWS Secrets Manager
// implementation when you're ready — the rest of the codebase doesn't care.
type SecretStore interface {
	Get(key string) (string, error)
}

// EnvSecretStore reads secrets from environment variables.
// Good enough for local dev; replace with AWS SM for prod use.
type EnvSecretStore struct{}

func NewEnvSecretStore() *EnvSecretStore {
	return &EnvSecretStore{}
}

func (s *EnvSecretStore) Get(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("secret %q not found in environment", key)
	}
	return val, nil
}
