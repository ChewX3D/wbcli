package cli

import (
	"strings"
	"testing"
)

func TestReadCredentialPairFromReader(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		maxBytes       int64
		expectedKey    string
		expectedSecret string
		expectedErr    error
	}{
		{
			name:           "valid payload",
			input:          "key-1\nsecret-1\n",
			maxBytes:       1024,
			expectedKey:    "key-1",
			expectedSecret: "secret-1",
		},
		{
			name:           "valid payload with crlf",
			input:          "key-1\r\nsecret-1\r\n",
			maxBytes:       1024,
			expectedKey:    "key-1",
			expectedSecret: "secret-1",
		},
		{
			name:        "empty payload",
			input:       "",
			maxBytes:    1024,
			expectedErr: ErrCredentialInputMissing,
		},
		{
			name:        "single line payload",
			input:       "key-only\n",
			maxBytes:    1024,
			expectedErr: ErrCredentialInputFormat,
		},
		{
			name:        "extra lines payload",
			input:       "key\nsecret\nextra\n",
			maxBytes:    1024,
			expectedErr: ErrCredentialInputFormat,
		},
		{
			name:        "oversized payload",
			input:       strings.Repeat("a", 20),
			maxBytes:    10,
			expectedErr: ErrCredentialInputTooLarge,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parsed, err := ReadCredentialPairFromReader(strings.NewReader(testCase.input), testCase.maxBytes)
			if testCase.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", testCase.expectedErr)
				}
				if err != testCase.expectedErr {
					t.Fatalf("expected error %v, got %v", testCase.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if parsed.APIKey != testCase.expectedKey {
				t.Fatalf("expected key %q, got %q", testCase.expectedKey, parsed.APIKey)
			}
			if string(parsed.APISecret) != testCase.expectedSecret {
				t.Fatalf("expected secret %q, got %q", testCase.expectedSecret, string(parsed.APISecret))
			}
		})
	}
}
