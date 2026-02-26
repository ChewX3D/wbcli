package auth

import "errors"

var (
	// ErrAPIKeyRequired indicates missing API key data.
	ErrAPIKeyRequired = errors.New("api key is required")
	// ErrAPISecretRequired indicates missing API secret data.
	ErrAPISecretRequired = errors.New("api secret is required")
)

// Credential contains auth material used for WhiteBIT requests.
type Credential struct {
	APIKey    string
	APISecret []byte
}

// Validate validates credential fields.
func (credential Credential) Validate() error {
	if credential.APIKey == "" {
		return ErrAPIKeyRequired
	}
	if len(credential.APISecret) == 0 {
		return ErrAPISecretRequired
	}

	return nil
}

// APIKeyHint returns a short safe representation of an API key.
func APIKeyHint(apiKey string) string {
	if apiKey == "" {
		return ""
	}
	if len(apiKey) <= 4 {
		return "****"
	}

	return apiKey[:2] + "***" + apiKey[len(apiKey)-2:]
}

// WipeBytes clears sensitive byte slices in place.
func WipeBytes(value []byte) {
	for index := range value {
		value[index] = 0
	}
}
