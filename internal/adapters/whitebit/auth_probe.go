package whitebit

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

const (
	defaultBaseURL         = "https://whitebit.com"
	hedgeModeProbePath     = "/api/v4/collateral-account/hedge-mode"
	defaultHTTPTimeout     = 10 * time.Second
	defaultProbeBodyMaxLen = 4096
)

// NonceSource generates strictly increasing nonces.
type NonceSource interface {
	Next() int64
}

// MonotonicUnixMilliNonceSource generates process-local monotonic unix-millisecond nonces.
type MonotonicUnixMilliNonceSource struct {
	mu   sync.Mutex
	last int64
}

// Next returns a strictly increasing nonce.
func (source *MonotonicUnixMilliNonceSource) Next() int64 {
	now := time.Now().UnixMilli()

	source.mu.Lock()
	defer source.mu.Unlock()

	if now <= source.last {
		now = source.last + 1
	}
	source.last = now

	return now
}

// AuthProbe validates WhiteBIT credentials with a signed hedge-mode request.
type AuthProbe struct {
	baseURL     string
	httpClient  *http.Client
	nonceSource NonceSource
}

type privateRequestPayload struct {
	Request string `json:"request"`
	Nonce   int64  `json:"nonce"`
}

// NewDefaultAuthProbe constructs AuthProbe with production defaults.
func NewDefaultAuthProbe() *AuthProbe {
	return NewAuthProbe(defaultBaseURL, nil, nil)
}

// NewAuthProbe constructs AuthProbe with dependency injection hooks for tests.
func NewAuthProbe(baseURL string, httpClient *http.Client, nonceSource NonceSource) *AuthProbe {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultBaseURL
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}
	if nonceSource == nil {
		nonceSource = &MonotonicUnixMilliNonceSource{}
	}

	return &AuthProbe{
		baseURL:     strings.TrimRight(baseURL, "/"),
		httpClient:  httpClient,
		nonceSource: nonceSource,
	}
}

// Probe executes a signed call to WhiteBIT hedge-mode endpoint.
func (probe *AuthProbe) Probe(ctx context.Context, credential domainauth.Credential) error {
	if err := credential.Validate(); err != nil {
		return err
	}

	requestBody := privateRequestPayload{
		Request: hedgeModeProbePath,
		Nonce:   probe.nonceSource.Next(),
	}
	rawBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("encode auth probe payload: %w", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(rawBody)
	signature := signPayload(encodedPayload, credential.APISecret)
	endpointURL := probe.baseURL + hedgeModeProbePath

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL, bytes.NewReader(rawBody))
	if err != nil {
		return fmt.Errorf("build auth probe request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-TXC-APIKEY", credential.APIKey)
	request.Header.Set("X-TXC-PAYLOAD", encodedPayload)
	request.Header.Set("X-TXC-SIGNATURE", signature)

	response, err := probe.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("%w: request failed", ports.ErrAuthProbeUnavailable)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, defaultProbeBodyMaxLen))
		_ = response.Body.Close()
	}()

	switch {
	case response.StatusCode >= 200 && response.StatusCode <= 299:
		return nil
	case response.StatusCode == http.StatusUnauthorized:
		return ports.ErrAuthProbeUnauthorized
	case response.StatusCode == http.StatusForbidden:
		return ports.ErrAuthProbeForbidden
	case response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= 500:
		return ports.ErrAuthProbeUnavailable
	default:
		return fmt.Errorf("auth probe unexpected status: %d", response.StatusCode)
	}
}

func signPayload(encodedPayload string, secret []byte) string {
	mac := hmac.New(sha512.New, secret)
	_, _ = mac.Write([]byte(encodedPayload))
	return hex.EncodeToString(mac.Sum(nil))
}
