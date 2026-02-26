package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
)

// CurrentResult returns the active profile metadata.
type CurrentResult struct {
	Profile    string
	Backend    string
	APIKeyHint string
	UpdatedAt  time.Time
}

// CurrentService provides active profile view.
type CurrentService struct {
	profileStore ports.ProfileStore
}

// NewCurrentService constructs CurrentService.
func NewCurrentService(profileStore ports.ProfileStore) *CurrentService {
	return &CurrentService{profileStore: profileStore}
}

// Execute returns current active profile metadata.
func (service *CurrentService) Execute(ctx context.Context) (CurrentResult, error) {
	activeProfile, err := service.profileStore.GetActiveProfile(ctx)
	if err != nil {
		return CurrentResult{}, fmt.Errorf("read active profile: %w", err)
	}
	if activeProfile == "" {
		return CurrentResult{}, ErrNoActiveProfile
	}

	profile, found, err := service.profileStore.GetProfile(ctx, activeProfile)
	if err != nil {
		return CurrentResult{}, fmt.Errorf("read active profile metadata: %w", err)
	}
	if !found {
		return CurrentResult{}, ErrProfileNotFound
	}

	return CurrentResult{
		Profile:    activeProfile,
		Backend:    profile.Backend,
		APIKeyHint: profile.APIKeyHint,
		UpdatedAt:  profile.UpdatedAt,
	}, nil
}
