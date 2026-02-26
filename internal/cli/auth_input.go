package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const defaultMaxCredentialPayloadBytes int64 = 16 * 1024

var (
	// ErrCredentialInputMissing indicates empty stdin payload.
	ErrCredentialInputMissing = errors.New("credential stdin payload is required")
	// ErrCredentialInputFormat indicates invalid stdin payload shape.
	ErrCredentialInputFormat = errors.New("credential stdin payload format is invalid")
	// ErrCredentialInputTooLarge indicates oversized stdin payload.
	ErrCredentialInputTooLarge = errors.New("credential stdin payload is too large")
)

// StdinCredentialInput holds parsed key/secret values.
type StdinCredentialInput struct {
	APIKey    string
	APISecret []byte
}

// ReadCredentialPairFromReader parses stdin payload as exactly two lines: key, secret.
func ReadCredentialPairFromReader(reader io.Reader, maxBytes int64) (StdinCredentialInput, error) {
	if maxBytes <= 0 {
		maxBytes = defaultMaxCredentialPayloadBytes
	}

	limitedReader := io.LimitReader(reader, maxBytes+1)
	payload, err := io.ReadAll(limitedReader)
	if err != nil {
		return StdinCredentialInput{}, fmt.Errorf("read credential stdin payload: %w", err)
	}
	if int64(len(payload)) > maxBytes {
		return StdinCredentialInput{}, ErrCredentialInputTooLarge
	}
	if len(payload) == 0 {
		return StdinCredentialInput{}, ErrCredentialInputMissing
	}

	normalizedPayload := normalizePayload(payload)
	if normalizedPayload == "" {
		return StdinCredentialInput{}, ErrCredentialInputMissing
	}

	lines := strings.Split(normalizedPayload, "\n")
	if len(lines) != 2 {
		return StdinCredentialInput{}, ErrCredentialInputFormat
	}
	apiKey := strings.TrimSuffix(lines[0], "\r")
	apiSecret := strings.TrimSuffix(lines[1], "\r")
	if apiKey == "" || apiSecret == "" {
		return StdinCredentialInput{}, ErrCredentialInputFormat
	}

	return StdinCredentialInput{
		APIKey:    apiKey,
		APISecret: []byte(apiSecret),
	}, nil
}

func normalizePayload(payload []byte) string {
	input := string(payload)
	if strings.HasSuffix(input, "\r\n") {
		input = strings.TrimSuffix(input, "\r\n")
	} else if strings.HasSuffix(input, "\n") {
		input = strings.TrimSuffix(input, "\n")
	}

	return input
}
