package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

type fakeCredentialStore struct {
	backendName string
	creds       map[string]domainauth.Credential
}

func (store *fakeCredentialStore) BackendName() string {
	return store.backendName
}

func (store *fakeCredentialStore) Save(_ context.Context, profile string, credential domainauth.Credential) error {
	if store.creds == nil {
		store.creds = map[string]domainauth.Credential{}
	}
	store.creds[profile] = domainauth.Credential{APIKey: credential.APIKey, APISecret: append([]byte(nil), credential.APISecret...)}
	return nil
}

func (store *fakeCredentialStore) Load(_ context.Context, profile string) (domainauth.Credential, error) {
	credential, ok := store.creds[profile]
	if !ok {
		return domainauth.Credential{}, ports.ErrCredentialNotFound
	}
	return domainauth.Credential{APIKey: credential.APIKey, APISecret: append([]byte(nil), credential.APISecret...)}, nil
}

func (store *fakeCredentialStore) Exists(_ context.Context, profile string) (bool, error) {
	_, ok := store.creds[profile]
	return ok, nil
}

func (store *fakeCredentialStore) Delete(_ context.Context, profile string) error {
	delete(store.creds, profile)
	return nil
}

type fakeProfileStore struct {
	profiles      map[string]ports.ProfileMetadata
	activeProfile string
}

func (store *fakeProfileStore) UpsertProfile(_ context.Context, profile ports.ProfileMetadata) error {
	if store.profiles == nil {
		store.profiles = map[string]ports.ProfileMetadata{}
	}
	store.profiles[profile.Name] = profile
	return nil
}

func (store *fakeProfileStore) GetProfile(_ context.Context, profile string) (ports.ProfileMetadata, bool, error) {
	metadata, ok := store.profiles[profile]
	return metadata, ok, nil
}

func (store *fakeProfileStore) ListProfiles(_ context.Context) ([]ports.ProfileMetadata, error) {
	profiles := make([]ports.ProfileMetadata, 0, len(store.profiles))
	for _, metadata := range store.profiles {
		profiles = append(profiles, metadata)
	}
	return profiles, nil
}

func (store *fakeProfileStore) DeleteProfile(_ context.Context, profile string) error {
	delete(store.profiles, profile)
	return nil
}

func (store *fakeProfileStore) SetActiveProfile(_ context.Context, profile string) error {
	store.activeProfile = profile
	return nil
}

func (store *fakeProfileStore) GetActiveProfile(_ context.Context) (string, error) {
	return store.activeProfile, nil
}

func (store *fakeProfileStore) ClearActiveProfile(_ context.Context) error {
	store.activeProfile = ""
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
	profileStore := &fakeProfileStore{}
	service := NewLoginService(credentialStore, profileStore, fixedClock{now: time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)})

	secret := []byte("secret-1")
	result, err := service.Execute(context.Background(), LoginRequest{
		Profile:   "prod",
		APIKey:    "api-key-1",
		APISecret: secret,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Profile != "prod" {
		t.Fatalf("expected profile prod, got %q", result.Profile)
	}
	if result.Backend != "os-keychain" {
		t.Fatalf("expected backend os-keychain, got %q", result.Backend)
	}
	if result.APIKeyHint == "" {
		t.Fatalf("expected api key hint")
	}

	storedProfile, ok := profileStore.profiles["prod"]
	if !ok {
		t.Fatalf("expected profile metadata to be stored")
	}
	if storedProfile.APIKeyHint == "" {
		t.Fatalf("expected api key hint in metadata")
	}
	if storedProfile.Name != "prod" {
		t.Fatalf("expected profile metadata name prod, got %q", storedProfile.Name)
	}
	if storedProfile.Backend != "os-keychain" {
		t.Fatalf("expected profile metadata backend os-keychain, got %q", storedProfile.Backend)
	}
	if got := string(credentialStore.creds["prod"].APISecret); got != "secret-1" {
		t.Fatalf("expected stored secret secret-1, got %q", got)
	}
	if profileStore.activeProfile != "prod" {
		t.Fatalf("expected active profile prod, got %q", profileStore.activeProfile)
	}
	if string(secret) != "\x00\x00\x00\x00\x00\x00\x00\x00" {
		t.Fatalf("expected request secret bytes to be wiped, got %q", string(secret))
	}
}

func TestLoginServiceExecuteRejectsOverwriteWithoutForce(t *testing.T) {
	credentialStore := &fakeCredentialStore{
		backendName: "os-keychain",
		creds: map[string]domainauth.Credential{
			"prod": {APIKey: "old", APISecret: []byte("old-secret")},
		},
	}
	profileStore := &fakeProfileStore{}
	service := NewLoginService(credentialStore, profileStore, fixedClock{now: time.Now()})

	_, err := service.Execute(context.Background(), LoginRequest{
		Profile:   "prod",
		APIKey:    "new",
		APISecret: []byte("new-secret"),
	})
	if !errors.Is(err, ErrCredentialAlreadyExists) {
		t.Fatalf("expected ErrCredentialAlreadyExists, got %v", err)
	}
}
