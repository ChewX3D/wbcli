package whitebit

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

const (
	defaultBaseURL      = "https://whitebit.com"
	defaultHTTPTimeout  = 10 * time.Second
	maxResponseBodySize = 64 * 1024
)

var (
	// ErrAPIAuth indicates authentication or authorization failure.
	ErrAPIAuth = errors.New("whitebit api auth error")
	// ErrUnauthorized indicates invalid API key/secret.
	ErrUnauthorized = errors.New("whitebit unauthorized")
	// ErrForbidden indicates valid credential without required permission.
	ErrForbidden = errors.New("whitebit forbidden")
	// ErrAPIValidation indicates request schema/validation failure.
	ErrAPIValidation = errors.New("whitebit api validation error")
	// ErrAPIBusinessRule indicates domain/business rejection by WhiteBIT.
	ErrAPIBusinessRule = errors.New("whitebit api business rule error")
	// ErrAPITransport indicates temporary transport/server/rate-limit failure.
	ErrAPITransport = errors.New("whitebit api transport error")
)

// HTTPDoer executes HTTP requests. It enables client testability.
type HTTPDoer interface {
	Do(request *http.Request) (*http.Response, error)
}

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

// Client executes signed private WhiteBIT HTTP API requests.
type Client struct {
	baseURL     string
	httpDoer    HTTPDoer
	nonceSource NonceSource
}

// NewDefaultClient constructs Client with production defaults.
func NewDefaultClient() *Client {
	return NewClient(defaultBaseURL, nil, nil)
}

// NewClient constructs Client with injectable dependencies for tests.
func NewClient(baseURL string, httpDoer HTTPDoer, nonceSource NonceSource) *Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultBaseURL
	}
	if httpDoer == nil {
		httpDoer = &http.Client{Timeout: defaultHTTPTimeout}
	}
	if nonceSource == nil {
		nonceSource = &MonotonicUnixMilliNonceSource{}
	}

	return &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		httpDoer:    httpDoer,
		nonceSource: nonceSource,
	}
}

type privateEnvelope struct {
	Request string `json:"request"`
	Nonce   string `json:"nonce"`
}

type whitebitErrorPayload struct {
	Message string `json:"message"`
}

func (client *Client) nextPrivateEnvelope(path string) privateEnvelope {
	return privateEnvelope{
		Request: path,
		Nonce:   strconv.FormatInt(client.nonceSource.Next(), 10),
	}
}

func (client *Client) doPrivateRequest(
	ctx context.Context,
	credential domainauth.Credential,
	path string,
	requestPayload any,
	responsePayload any,
) error {
	if err := credential.Validate(); err != nil {
		return err
	}

	rawBody, err := json.Marshal(requestPayload)
	if err != nil {
		return fmt.Errorf("encode request body: %w", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(rawBody)
	signature := signPayload(encodedPayload, credential.APISecret)

	endpointURL := client.baseURL + path
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL, bytes.NewReader(rawBody))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-TXC-APIKEY", credential.APIKey)
	request.Header.Set("X-TXC-PAYLOAD", encodedPayload)
	request.Header.Set("X-TXC-SIGNATURE", signature)

	response, err := client.httpDoer.Do(request)
	if err != nil {
		return fmt.Errorf("%w: request failed", ErrAPITransport)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, maxResponseBodySize))
		_ = response.Body.Close()
	}()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBodySize))
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		return mapHTTPStatusError(response.StatusCode, responseBody)
	}

	if responsePayload == nil || len(responseBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(responseBody, responsePayload); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}

func mapHTTPStatusError(statusCode int, body []byte) error {
	responseMessage := extractErrorMessage(body)
	wrapStatus := func(base error) error {
		if responseMessage == "" {
			return fmt.Errorf("%w: status %d", base, statusCode)
		}
		return fmt.Errorf("%w: status %d: %s", base, statusCode, responseMessage)
	}

	switch {
	case statusCode == http.StatusUnauthorized:
		return fmt.Errorf("%w: %w", wrapStatus(ErrAPIAuth), ErrUnauthorized)
	case statusCode == http.StatusForbidden:
		return fmt.Errorf("%w: %w", wrapStatus(ErrAPIAuth), ErrForbidden)
	case statusCode == http.StatusBadRequest || statusCode == http.StatusUnprocessableEntity:
		return wrapStatus(ErrAPIValidation)
	case statusCode == http.StatusTooManyRequests || statusCode >= 500:
		return wrapStatus(ErrAPITransport)
	case statusCode >= 400 && statusCode <= 499:
		return wrapStatus(ErrAPIBusinessRule)
	default:
		return wrapStatus(ErrAPITransport)
	}
}

func extractErrorMessage(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var payload whitebitErrorPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}

	return strings.TrimSpace(payload.Message)
}

func signPayload(encodedPayload string, secret []byte) string {
	mac := hmac.New(sha512.New, secret)
	_, _ = mac.Write([]byte(encodedPayload))
	return hex.EncodeToString(mac.Sum(nil))
}
