package whitebit

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

type fixedNonceSource struct {
	value int64
}

func (source fixedNonceSource) Next() int64 {
	return source.value
}

func TestAuthProbeProbeSuccessSignsRequest(t *testing.T) {
	credential := domainauth.Credential{
		APIKey:    "test-key",
		APISecret: []byte("test-secret"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", request.Method)
		}
		if request.URL.Path != hedgeModeProbePath {
			t.Fatalf("expected path %s, got %s", hedgeModeProbePath, request.URL.Path)
		}
		if got := request.Header.Get("X-TXC-APIKEY"); got != credential.APIKey {
			t.Fatalf("expected api key header %q, got %q", credential.APIKey, got)
		}

		bodyBytes, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		payloadHeader := request.Header.Get("X-TXC-PAYLOAD")
		if payloadHeader == "" {
			t.Fatalf("expected payload header")
		}
		decodedPayload, err := base64.StdEncoding.DecodeString(payloadHeader)
		if err != nil {
			t.Fatalf("decode payload header: %v", err)
		}
		if string(decodedPayload) != string(bodyBytes) {
			t.Fatalf("expected body to match decoded payload header")
		}

		expectedSignature := signPayload(payloadHeader, credential.APISecret)
		if got := request.Header.Get("X-TXC-SIGNATURE"); got != expectedSignature {
			t.Fatalf("expected signature %q, got %q", expectedSignature, got)
		}

		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	probe := NewAuthProbe(server.URL, server.Client(), fixedNonceSource{value: 1700000000000})
	if err := probe.Probe(context.Background(), credential); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAuthProbeProbeMapsUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	probe := NewAuthProbe(server.URL, server.Client(), fixedNonceSource{value: 1})
	err := probe.Probe(context.Background(), domainauth.Credential{
		APIKey:    "bad",
		APISecret: []byte("bad-secret"),
	})
	if !errors.Is(err, ports.ErrAuthProbeUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

func TestAuthProbeProbeMapsForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	probe := NewAuthProbe(server.URL, server.Client(), fixedNonceSource{value: 1})
	err := probe.Probe(context.Background(), domainauth.Credential{
		APIKey:    "no-perms",
		APISecret: []byte("secret"),
	})
	if !errors.Is(err, ports.ErrAuthProbeForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}

func TestAuthProbeProbeMapsUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	probe := NewAuthProbe(server.URL, server.Client(), fixedNonceSource{value: 1})
	err := probe.Probe(context.Background(), domainauth.Credential{
		APIKey:    "key",
		APISecret: []byte("secret"),
	})
	if !errors.Is(err, ports.ErrAuthProbeUnavailable) {
		t.Fatalf("expected unavailable error, got %v", err)
	}
}
