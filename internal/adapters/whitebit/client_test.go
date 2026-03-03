package whitebit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func boolPtr(value bool) *bool {
	return &value
}

func TestClientGetCollateralAccountHedgeMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != collateralAccountHedgeModePath {
			t.Fatalf("expected path %s, got %s", collateralAccountHedgeModePath, request.URL.Path)
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"hedgeMode":true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client(), fixedNonceSource{value: 1})
	response, err := client.GetCollateralAccountHedgeMode(context.Background(), domainauth.Credential{
		APIKey:    "public-key",
		APISecret: []byte("secret-key"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !response.HedgeMode {
		t.Fatalf("expected hedge mode true")
	}
}

func TestClientStatusErrorMapping(t *testing.T) {
	testCases := []struct {
		name           string
		statusCode     int
		expectedErrors []error
	}{
		{
			name:           "unauthorized",
			statusCode:     http.StatusUnauthorized,
			expectedErrors: []error{ErrAPIAuth, ErrUnauthorized},
		},
		{
			name:           "forbidden",
			statusCode:     http.StatusForbidden,
			expectedErrors: []error{ErrAPIAuth, ErrForbidden},
		},
		{
			name:           "validation",
			statusCode:     http.StatusUnprocessableEntity,
			expectedErrors: []error{ErrAPIValidation},
		},
		{
			name:           "business rule",
			statusCode:     http.StatusConflict,
			expectedErrors: []error{ErrAPIBusinessRule},
		},
		{
			name:           "transport",
			statusCode:     http.StatusServiceUnavailable,
			expectedErrors: []error{ErrAPITransport},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(testCase.statusCode)
				_, _ = writer.Write([]byte(`{"message":"failed"}`))
			}))
			defer server.Close()

			client := NewClient(server.URL, server.Client(), fixedNonceSource{value: 1})
			_, err := client.GetCollateralAccountHedgeMode(context.Background(), domainauth.Credential{
				APIKey:    "public-key",
				APISecret: []byte("secret-key"),
			})
			if err == nil {
				t.Fatalf("expected error")
			}

			for _, expectedErr := range testCase.expectedErrors {
				if !errors.Is(err, expectedErr) {
					t.Fatalf("expected error %v, got %v", expectedErr, err)
				}
			}
		})
	}
}

func TestMapHTTPStatusErrorValidationIncludesFieldDetails(t *testing.T) {
	err := mapHTTPStatusError(http.StatusUnprocessableEntity, []byte(`{
		"code": 32,
		"message": "Validation failed",
		"errors": {
			"market": ["The selected market is invalid."],
			"amount": ["The amount field is required."]
		}
	}`))
	if !errors.Is(err, ErrAPIValidation) {
		t.Fatalf("expected validation error category, got %v", err)
	}
	if !strings.Contains(err.Error(), "code 32:") {
		t.Fatalf("expected code prefix, got %v", err)
	}
	if !strings.Contains(err.Error(), "Validation failed") {
		t.Fatalf("expected generic message, got %v", err)
	}
	if !strings.Contains(err.Error(), "market: The selected market is invalid.") {
		t.Fatalf("expected market detail, got %v", err)
	}
	if !strings.Contains(err.Error(), "amount: The amount field is required.") {
		t.Fatalf("expected amount detail, got %v", err)
	}
}

func TestExtractErrorMessageWithCodeField(t *testing.T) {
	message := extractErrorMessage([]byte(`{
		"code": 37,
		"message": "Validation failed",
		"errors": {"ioc": ["ioc cannot be combined with postOnly"]}
	}`))
	if !strings.HasPrefix(message, "code 37: ") {
		t.Fatalf("expected code prefix, got %q", message)
	}
	if !strings.Contains(message, "Validation failed") {
		t.Fatalf("expected message body, got %q", message)
	}
}

func TestExtractErrorMessageWithoutCodeField(t *testing.T) {
	message := extractErrorMessage([]byte(`{"message": "Something failed"}`))
	if strings.Contains(message, "code") {
		t.Fatalf("expected no code prefix, got %q", message)
	}
	if message != "Something failed" {
		t.Fatalf("expected plain message, got %q", message)
	}
}

func TestExtractErrorMessageFromErrorsOnlyPayload(t *testing.T) {
	message := extractErrorMessage([]byte(`{
		"errors": {
			"price": ["The price field is required."]
		}
	}`))
	if !strings.Contains(message, "price: The price field is required.") {
		t.Fatalf("expected errors-only detail, got %q", message)
	}
}

func TestCredentialVerifierAdapterVerifyMapsErrors(t *testing.T) {
	testCases := []struct {
		name         string
		statusCode   int
		expectedCode ports.ErrorCode
	}{
		{
			name:         "unauthorized",
			statusCode:   http.StatusUnauthorized,
			expectedCode: ports.CodeUnauthorized,
		},
		{
			name:         "forbidden",
			statusCode:   http.StatusForbidden,
			expectedCode: ports.CodeForbidden,
		},
		{
			name:         "unavailable",
			statusCode:   http.StatusServiceUnavailable,
			expectedCode: ports.CodeUnavailable,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(testCase.statusCode)
			}))
			defer server.Close()

			client := NewClient(server.URL, server.Client(), fixedNonceSource{value: 1})
			adapter := NewCredentialVerifierAdapter(client)
			_, err := adapter.Verify(context.Background(), domainauth.Credential{
				APIKey:    "public-key",
				APISecret: []byte("secret-key"),
			})
			var apiErr *ports.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T: %v", err, err)
			}
			if apiErr.Code != testCase.expectedCode {
				t.Fatalf("expected code %v, got %v", testCase.expectedCode, apiErr.Code)
			}
		})
	}
}

func TestCredentialVerifierAdapterVerifyReturnsEndpointOnSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"hedgeMode":true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client(), fixedNonceSource{value: 1})
	adapter := NewCredentialVerifierAdapter(client)

	result, err := adapter.Verify(context.Background(), domainauth.Credential{
		APIKey:    "public-key",
		APISecret: []byte("secret-key"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Endpoint != collateralAccountHedgeModePath {
		t.Fatalf("expected endpoint %q, got %q", collateralAccountHedgeModePath, result.Endpoint)
	}
	if !result.HedgeMode {
		t.Fatalf("expected hedge_mode=true, got false")
	}
}

func TestCredentialVerifierAdapterUnauthorizedActionDeniedMapsToInsufficientAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(`{"message":"This API Key is not authorized to perform this action."}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client(), fixedNonceSource{value: 1})
	adapter := NewCredentialVerifierAdapter(client)

	_, err := adapter.Verify(context.Background(), domainauth.Credential{
		APIKey:    "public-key",
		APISecret: []byte("secret-key"),
	})
	if err == nil {
		t.Fatalf("expected error")
	}

	var apiErr *ports.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.Code != ports.CodeForbidden {
		t.Fatalf("expected CodeForbidden, got %s", apiErr.Code)
	}
	if !strings.Contains(apiErr.Details, collateralAccountHedgeModePath) {
		t.Fatalf("expected endpoint %q in details, got %q", collateralAccountHedgeModePath, apiErr.Details)
	}
	if !strings.Contains(strings.ToLower(apiErr.Details), "not authorized to perform this action") {
		t.Fatalf("expected action denied detail, got %q", apiErr.Details)
	}
}

func TestClientValidatesEnumAndOrderRules(t *testing.T) {
	client := NewClient("https://whitebit.com", &http.Client{}, fixedNonceSource{value: 1})
	credential := domainauth.Credential{
		APIKey:    "public-key",
		APISecret: []byte("secret-key"),
	}

	_, err := client.PlaceCollateralLimitOrder(context.Background(), credential, CollateralLimitOrderRequest{
		Market: "BTC_PERP",
		Side:   OrderSide("hold"),
		Amount: "0.001",
		Price:  "50000",
	})
	if !errors.Is(err, ErrInvalidOrderSide) {
		t.Fatalf("expected invalid side error, got %v", err)
	}

	_, err = client.PlaceCollateralLimitOrder(context.Background(), credential, CollateralLimitOrderRequest{
		Market:       "BTC_PERP",
		Side:         OrderSideBuy,
		Amount:       "0.001",
		Price:        "50000",
		PositionSide: PositionSide("hedge"),
	})
	if !errors.Is(err, ErrInvalidPositionSide) {
		t.Fatalf("expected invalid position side error, got %v", err)
	}

	_, err = client.PlaceCollateralLimitOrder(context.Background(), credential, CollateralLimitOrderRequest{
		Market:   "BTC_PERP",
		Side:     OrderSideBuy,
		Amount:   "0.001",
		Price:    "50000",
		PostOnly: boolPtr(true),
		IOC:      boolPtr(true),
	})
	if !errors.Is(err, ErrIOCConflict) {
		t.Fatalf("expected ioc conflict error for postOnly+ioc, got %v", err)
	}

	_, err = client.PlaceCollateralLimitOrder(context.Background(), credential, CollateralLimitOrderRequest{
		Market: "BTC_PERP",
		Side:   OrderSideBuy,
		Amount: "0.001",
		Price:  "50000",
		IOC:    boolPtr(true),
		RPI:    boolPtr(true),
	})
	if !errors.Is(err, ErrIOCConflict) {
		t.Fatalf("expected ioc conflict error for ioc+rpi, got %v", err)
	}
}
