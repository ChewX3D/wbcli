package whitebit

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

func boolPtr(value bool) *bool {
	return &value
}

func decodePayloadHeader(t *testing.T, request *http.Request) map[string]any {
	t.Helper()

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("read request body: %v", err)
	}
	payloadHeader := request.Header.Get("X-TXC-PAYLOAD")
	if payloadHeader == "" {
		t.Fatalf("missing payload header")
	}

	decodedPayload, err := base64.StdEncoding.DecodeString(payloadHeader)
	if err != nil {
		t.Fatalf("decode payload header: %v", err)
	}
	if string(decodedPayload) != string(bodyBytes) {
		t.Fatalf("expected request body to match decoded payload header")
	}

	var payloadMap map[string]any
	if err := json.Unmarshal(decodedPayload, &payloadMap); err != nil {
		t.Fatalf("decode payload json: %v", err)
	}

	return payloadMap
}

func TestClientPlaceCollateralLimitOrderSignsHeadersAndBody(t *testing.T) {
	credential := domainauth.Credential{
		APIKey:    "public-key",
		APISecret: []byte("secret-key"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", request.Method)
		}
		if request.URL.Path != collateralLimitOrderPath {
			t.Fatalf("expected path %s, got %s", collateralLimitOrderPath, request.URL.Path)
		}
		if got := request.Header.Get("X-TXC-APIKEY"); got != credential.APIKey {
			t.Fatalf("expected api key header %q, got %q", credential.APIKey, got)
		}
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("expected content-type application/json, got %q", got)
		}

		payload := decodePayloadHeader(t, request)
		if payload["request"] != collateralLimitOrderPath {
			t.Fatalf("expected request field %q, got %#v", collateralLimitOrderPath, payload["request"])
		}
		if payload["nonce"] != "1700000000000" {
			t.Fatalf("expected nonce 1700000000000, got %#v", payload["nonce"])
		}
		if payload["side"] != string(OrderSideBuy) {
			t.Fatalf("expected side buy, got %#v", payload["side"])
		}
		if payload["positionSide"] != string(PositionSideLong) {
			t.Fatalf("expected position side long, got %#v", payload["positionSide"])
		}

		payloadHeader := request.Header.Get("X-TXC-PAYLOAD")
		expectedSignature := signPayload(payloadHeader, credential.APISecret)
		if got := request.Header.Get("X-TXC-SIGNATURE"); got != expectedSignature {
			t.Fatalf("expected signature %q, got %q", expectedSignature, got)
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"orderId":12345}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client(), fixedNonceSource{value: 1700000000000})
	response, err := client.PlaceCollateralLimitOrder(context.Background(), credential, CollateralLimitOrderRequest{
		Market:       "BTC_PERP",
		Side:         OrderSideBuy,
		Amount:       "0.001",
		Price:        "50000",
		PositionSide: PositionSideLong,
		PostOnly:     boolPtr(true),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(response) != `{"orderId":12345}` {
		t.Fatalf("expected raw response payload, got %s", string(response))
	}
}

func TestClientPlaceCollateralBulkLimitOrderSendsOrdersPayload(t *testing.T) {
	credential := domainauth.Credential{
		APIKey:    "public-key",
		APISecret: []byte("secret-key"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != collateralBulkLimitOrderPath {
			t.Fatalf("expected path %s, got %s", collateralBulkLimitOrderPath, request.URL.Path)
		}

		payload := decodePayloadHeader(t, request)
		if payload["request"] != collateralBulkLimitOrderPath {
			t.Fatalf("expected request field %q, got %#v", collateralBulkLimitOrderPath, payload["request"])
		}
		if payload["nonce"] != "1700000000123" {
			t.Fatalf("expected nonce 1700000000123, got %#v", payload["nonce"])
		}

		orders, ok := payload["orders"].([]any)
		if !ok || len(orders) != 2 {
			t.Fatalf("expected 2 orders, got %#v", payload["orders"])
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client(), fixedNonceSource{value: 1700000000123})
	response, err := client.PlaceCollateralBulkLimitOrder(context.Background(), credential, CollateralBulkLimitOrderRequest{
		StopOnFail: boolPtr(true),
		Orders: []CollateralLimitOrderRequest{
			{
				Market:       "BTC_PERP",
				Side:         OrderSideBuy,
				Amount:       "0.001",
				Price:        "50000",
				PositionSide: PositionSideLong,
			},
			{
				Market:       "BTC_PERP",
				Side:         OrderSideSell,
				Amount:       "0.001",
				Price:        "50100",
				PositionSide: PositionSideShort,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(response) != `{"success":true}` {
		t.Fatalf("expected raw response payload, got %s", string(response))
	}
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

func TestClientProbeMapsErrorsForAuthLogin(t *testing.T) {
	testCases := []struct {
		name        string
		statusCode  int
		expectedErr error
	}{
		{
			name:        "unauthorized",
			statusCode:  http.StatusUnauthorized,
			expectedErr: ports.ErrAuthProbeUnauthorized,
		},
		{
			name:        "forbidden",
			statusCode:  http.StatusForbidden,
			expectedErr: ports.ErrAuthProbeForbidden,
		},
		{
			name:        "unavailable",
			statusCode:  http.StatusServiceUnavailable,
			expectedErr: ports.ErrAuthProbeUnavailable,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(testCase.statusCode)
			}))
			defer server.Close()

			client := NewClient(server.URL, server.Client(), fixedNonceSource{value: 1})
			err := client.Probe(context.Background(), domainauth.Credential{
				APIKey:    "public-key",
				APISecret: []byte("secret-key"),
			})
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("expected error %v, got %v", testCase.expectedErr, err)
			}
		})
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
	if !errors.Is(err, ErrPostOnlyIOCConflict) {
		t.Fatalf("expected postOnly/ioc conflict error, got %v", err)
	}
}
