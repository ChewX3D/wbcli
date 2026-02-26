package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

// LogoutRequest defines profile target for auth logout.
type LogoutRequest struct {
	Profile string
}

// LogoutResult contains safe logout result.
type LogoutResult struct {
	Profile string
}

// LogoutService deletes profile credentials.
type LogoutService struct {
	credentialStore ports.CredentialStore
	profileStore    ports.ProfileStore
}

// NewLogoutService constructs LogoutService.
func NewLogoutService(credentialStore ports.CredentialStore, profileStore ports.ProfileStore) *LogoutService {
	return &LogoutService{
		credentialStore: credentialStore,
		profileStore:    profileStore,
	}
}

// Execute removes profile credentials and metadata.
func (service *LogoutService) Execute(ctx context.Context, request LogoutRequest) (LogoutResult, error) {
	if err := domainauth.ValidateProfileName(request.Profile); err != nil {
		return LogoutResult{}, err
	}

	if err := service.credentialStore.Delete(ctx, request.Profile); err != nil {
		if !errors.Is(err, ports.ErrCredentialNotFound) {
			return LogoutResult{}, fmt.Errorf("delete credential: %w", err)
		}
	}
	if err := service.profileStore.DeleteProfile(ctx, request.Profile); err != nil {
		return LogoutResult{}, fmt.Errorf("delete profile metadata: %w", err)
	}

	activeProfile, err := service.profileStore.GetActiveProfile(ctx)
	if err != nil {
		return LogoutResult{}, fmt.Errorf("read active profile: %w", err)
	}
	if activeProfile == request.Profile {
		if err := service.profileStore.ClearActiveProfile(ctx); err != nil {
			return LogoutResult{}, fmt.Errorf("clear active profile: %w", err)
		}
	}

	return LogoutResult{Profile: request.Profile}, nil
}
