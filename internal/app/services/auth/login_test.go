package auth

import (
	"context"
	"testing"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

type fakeCredentialStore struct {
	backendName string
	credential  *domainauth.Credential
}

func (store *fakeCredentialStore) BackendName() string {
	return store.backendName
}

func (store *fakeCredentialStore) Save(_ context.Context, credential domainauth.Credential) error {
	copied := domainauth.Credential{APIKey: credential.APIKey, APISecret: append([]byte(nil), credential.APISecret...)}
	store.credential = &copied
	return nil
}

func (store *fakeCredentialStore) Load(_ context.Context) (domainauth.Credential, error) {
	if store.credential == nil {
		return domainauth.Credential{}, ports.ErrCredentialNotFound
	}
	return domainauth.Credential{APIKey: store.credential.APIKey, APISecret: append([]byte(nil), store.credential.APISecret...)}, nil
}

func (store *fakeCredentialStore) Exists(_ context.Context) (bool, error) {
	return store.credential != nil, nil
}

func (store *fakeCredentialStore) Delete(_ context.Context) error {
	store.credential = nil
	return nil
}

type fakeSessionStore struct {
	session *ports.SessionMetadata
}

func (store *fakeSessionStore) SaveSession(_ context.Context, session ports.SessionMetadata) error {
	copied := session
	store.session = &copied
	return nil
}

func (store *fakeSessionStore) GetSession(_ context.Context) (ports.SessionMetadata, bool, error) {
	if store.session == nil {
		return ports.SessionMetadata{}, false, nil
	}
	return *store.session, true, nil
}

func (store *fakeSessionStore) ClearSession(_ context.Context) error {
	store.session = nil
	return nil
}

type fixedClock struct {
	now time.Time
}

func (clock fixedClock) Now() time.Time {
	return clock.now
}

func TestLoginServiceExecuteSuccess(t *testing.T) {
	credentialStore := &fakeCredentialStore{backendName: "os-keychain"}
	sessionStore := &fakeSessionStore{}
	service := NewLoginService(credentialStore, sessionStore, fixedClock{now: time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)})

	secret := []byte("secret-1")
	result, err := service.Execute(context.Background(), LoginRequest{
		APIKey:    "api-key-1",
		APISecret: secret,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Backend != "os-keychain" {
		t.Fatalf("expected backend os-keychain, got %q", result.Backend)
	}
	if result.APIKeyHint == "" {
		t.Fatalf("expected api key hint")
	}
	if credentialStore.credential == nil {
		t.Fatalf("expected credential to be stored")
	}
	if got := string(credentialStore.credential.APISecret); got != "secret-1" {
		t.Fatalf("expected stored secret secret-1, got %q", got)
	}
	if sessionStore.session == nil {
		t.Fatalf("expected session metadata to be stored")
	}
	if sessionStore.session.Backend != "os-keychain" {
		t.Fatalf("expected session backend os-keychain, got %q", sessionStore.session.Backend)
	}
	if string(secret) != "\x00\x00\x00\x00\x00\x00\x00\x00" {
		t.Fatalf("expected request secret bytes to be wiped, got %q", string(secret))
	}
}

func TestLoginServiceExecuteOverwritesExistingCredential(t *testing.T) {
	credentialStore := &fakeCredentialStore{
		backendName: "os-keychain",
		credential:  &domainauth.Credential{APIKey: "old", APISecret: []byte("old-secret")},
	}
	sessionStore := &fakeSessionStore{}
	service := NewLoginService(credentialStore, sessionStore, fixedClock{now: time.Now()})

	result, err := service.Execute(context.Background(), LoginRequest{
		APIKey:    "new",
		APISecret: []byte("new-secret"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.APIKeyHint == "" {
		t.Fatalf("expected api key hint")
	}
	if credentialStore.credential == nil {
		t.Fatalf("expected credential to be stored")
	}
	if credentialStore.credential.APIKey != "new" {
		t.Fatalf("expected API key to be overwritten, got %q", credentialStore.credential.APIKey)
	}
}
