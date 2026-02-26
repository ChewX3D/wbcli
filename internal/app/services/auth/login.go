package auth

import (
	"context"
	"fmt"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

// LoginRequest is input for auth login use-case.
type LoginRequest struct {
	Profile   string
	APIKey    string
	APISecret []byte
	Force     bool
}

// LoginResult is safe output for auth login use-case.
type LoginResult struct {
	Profile    string
	Backend    string
	APIKeyHint string
	SavedAt    string
}

// LoginService stores profile credentials securely.
type LoginService struct {
	credentialStore ports.CredentialStore
	profileStore    ports.ProfileStore
	clock           ports.Clock
}

// NewLoginService constructs LoginService.
func NewLoginService(credentialStore ports.CredentialStore, profileStore ports.ProfileStore, clock ports.Clock) *LoginService {
	return &LoginService{
		credentialStore: credentialStore,
		profileStore:    profileStore,
		clock:           clock,
	}
}

// Execute validates login input and writes auth material into stores.
func (service *LoginService) Execute(ctx context.Context, request LoginRequest) (LoginResult, error) {
	if err := domainauth.ValidateProfileName(request.Profile); err != nil {
		return LoginResult{}, err
	}

	credential := domainauth.Credential{
		APIKey:    request.APIKey,
		APISecret: request.APISecret,
	}
	if err := credential.Validate(); err != nil {
		return LoginResult{}, err
	}
	defer domainauth.WipeBytes(request.APISecret)

	exists, err := service.credentialStore.Exists(ctx, request.Profile)
	if err != nil {
		return LoginResult{}, fmt.Errorf("check existing credential: %w", err)
	}
	if exists && !request.Force {
		return LoginResult{}, ErrCredentialAlreadyExists
	}

	previous, found, err := service.profileStore.GetProfile(ctx, request.Profile)
	if err != nil {
		return LoginResult{}, fmt.Errorf("read profile metadata: %w", err)
	}

	now := service.clock.Now().UTC()
	metadata := ports.ProfileMetadata{
		Name:       request.Profile,
		Backend:    service.credentialStore.BackendName(),
		APIKeyHint: domainauth.APIKeyHint(request.APIKey),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if found {
		metadata.CreatedAt = previous.CreatedAt
	}

	if err := service.credentialStore.Save(ctx, request.Profile, credential); err != nil {
		return LoginResult{}, fmt.Errorf("save credential: %w", err)
	}
	if err := service.profileStore.UpsertProfile(ctx, metadata); err != nil {
		return LoginResult{}, fmt.Errorf("save profile metadata: %w", err)
	}
	if err := service.profileStore.SetActiveProfile(ctx, request.Profile); err != nil {
		return LoginResult{}, fmt.Errorf("set active profile: %w", err)
	}

	return LoginResult{
		Profile:    request.Profile,
		Backend:    metadata.Backend,
		APIKeyHint: metadata.APIKeyHint,
		SavedAt:    now.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
