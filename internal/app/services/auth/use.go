package auth

import (
	"context"
	"fmt"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

// UseRequest defines input for selecting active profile.
type UseRequest struct {
	Profile string
}

// UseResult defines safe output for profile selection.
type UseResult struct {
	Profile string
}

// UseService sets active auth profile.
type UseService struct {
	credentialStore ports.CredentialStore
	profileStore    ports.ProfileStore
	clock           ports.Clock
}

// NewUseService constructs UseService.
func NewUseService(credentialStore ports.CredentialStore, profileStore ports.ProfileStore, clock ports.Clock) *UseService {
	return &UseService{
		credentialStore: credentialStore,
		profileStore:    profileStore,
		clock:           clock,
	}
}

// Execute sets active profile after metadata and credential checks.
func (service *UseService) Execute(ctx context.Context, request UseRequest) (UseResult, error) {
	if err := domainauth.ValidateProfileName(request.Profile); err != nil {
		return UseResult{}, err
	}

	metadata, found, err := service.profileStore.GetProfile(ctx, request.Profile)
	if err != nil {
		return UseResult{}, fmt.Errorf("read profile metadata: %w", err)
	}
	if !found {
		return UseResult{}, ErrProfileNotFound
	}

	exists, err := service.credentialStore.Exists(ctx, request.Profile)
	if err != nil {
		return UseResult{}, fmt.Errorf("check credential existence: %w", err)
	}
	if !exists {
		return UseResult{}, ports.ErrCredentialNotFound
	}

	now := service.clock.Now().UTC()
	metadata.UpdatedAt = now
	metadata.LastUsedAt = &now
	if err := service.profileStore.UpsertProfile(ctx, metadata); err != nil {
		return UseResult{}, fmt.Errorf("update profile metadata: %w", err)
	}

	if err := service.profileStore.SetActiveProfile(ctx, request.Profile); err != nil {
		return UseResult{}, fmt.Errorf("set active profile: %w", err)
	}

	return UseResult{Profile: request.Profile}, nil
}
