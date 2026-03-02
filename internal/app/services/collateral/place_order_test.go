package collateral

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

type fakeClock struct {
	now time.Time
}

func (clock fakeClock) Now() time.Time {
	return clock.now
}

type fakeCredentialStore struct {
	loadCredential domainauth.Credential
	loadErr        error
}

func (store *fakeCredentialStore) BackendName() string {
	return "os-keychain"
}

func (store *fakeCredentialStore) Save(context.Context, domainauth.Credential) error {
	return nil
}

func (store *fakeCredentialStore) Load(context.Context) (domainauth.Credential, error) {
	if store.loadErr != nil {
		return domainauth.Credential{}, store.loadErr
	}

	return store.loadCredential, nil
}

func (store *fakeCredentialStore) Exists(context.Context) (bool, error) {
	return true, nil
}

func (store *fakeCredentialStore) Delete(context.Context) error {
	return nil
}

type fakeOrderExecutor struct {
	lastCredential domainauth.Credential
	lastRequest    ports.CollateralLimitOrderRequest
	err            error
}

func (executor *fakeOrderExecutor) PlaceCollateralLimitOrder(
	_ context.Context,
	credential domainauth.Credential,
	request ports.CollateralLimitOrderRequest,
) (json.RawMessage, error) {
	executor.lastCredential = domainauth.Credential{
		APIKey:    credential.APIKey,
		APISecret: append([]byte(nil), credential.APISecret...),
	}
	executor.lastRequest = request

	if executor.err != nil {
		return nil, executor.err
	}

	return json.RawMessage(`{"status":"ok"}`), nil
}

func TestPlaceOrderServiceExecuteSuccess(t *testing.T) {
	credentialStore := &fakeCredentialStore{
		loadCredential: domainauth.Credential{
			APIKey:    "public-key",
			APISecret: []byte("secret-key"),
		},
	}
	orderExecutor := &fakeOrderExecutor{}
	service := NewPlaceOrderService(
		credentialStore,
		orderExecutor,
		fakeClock{now: time.Date(2026, 3, 2, 10, 0, 0, 123, time.UTC)},
	)

	result, err := service.Execute(context.Background(), PlaceOrderRequest{
		Market:        "BTC_PERP",
		Side:          "buy",
		Amount:        "0.01",
		Price:         "50000",
		ClientOrderID: "client-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Mode != "single" {
		t.Fatalf("expected mode single, got %q", result.Mode)
	}
	if result.OrdersPlanned != 1 || result.OrdersSubmitted != 1 || result.OrdersFailed != 0 {
		t.Fatalf("unexpected counters: %+v", result)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected empty errors, got %v", result.Errors)
	}
	if orderExecutor.lastRequest.PostOnly != true {
		t.Fatalf("expected postOnly=true")
	}
	if orderExecutor.lastRequest.ClientOrderID != "client-001" {
		t.Fatalf("expected client order id to pass through, got %q", orderExecutor.lastRequest.ClientOrderID)
	}
}

func TestPlaceOrderServiceExecuteCredentialLoadFailure(t *testing.T) {
	credentialStore := &fakeCredentialStore{loadErr: ports.ErrCredentialNotFound}
	orderExecutor := &fakeOrderExecutor{}
	service := NewPlaceOrderService(credentialStore, orderExecutor, fakeClock{now: time.Now()})

	_, err := service.Execute(context.Background(), PlaceOrderRequest{
		Market: "BTC_PERP",
		Side:   "buy",
		Amount: "0.01",
		Price:  "50000",
	})
	if !errors.Is(err, ports.ErrCredentialNotFound) {
		t.Fatalf("expected credential not found, got %v", err)
	}
}

func TestPlaceOrderServiceExecuteExecutorFailure(t *testing.T) {
	credentialStore := &fakeCredentialStore{
		loadCredential: domainauth.Credential{
			APIKey:    "public-key",
			APISecret: []byte("secret-key"),
		},
	}
	orderExecutor := &fakeOrderExecutor{err: errors.New("exchange rejected request")}
	service := NewPlaceOrderService(credentialStore, orderExecutor, fakeClock{now: time.Now()})

	_, err := service.Execute(context.Background(), PlaceOrderRequest{
		Market: "BTC_PERP",
		Side:   "buy",
		Amount: "0.01",
		Price:  "50000",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if got := err.Error(); got == "" {
		t.Fatalf("expected non-empty error")
	}
}
