package auth

import (
	"context"
	"fmt"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

// TestRequest defines profile target for auth test.
type TestRequest struct {
	Profile string
}

// TestService validates credentials against WhiteBIT via probe port.
type TestService struct {
	credentialStore ports.CredentialStore
	profileStore    ports.ProfileStore
	probe           ports.AuthProbe
}

// NewTestService constructs TestService.
func NewTestService(credentialStore ports.CredentialStore, profileStore ports.ProfileStore, probe ports.AuthProbe) *TestService {
	return &TestService{
		credentialStore: credentialStore,
		profileStore:    profileStore,
		probe:           probe,
	}
}

// Execute runs auth test when probe is available.
func (service *TestService) Execute(ctx context.Context, request TestRequest) error {
	if service.probe == nil {
		return ErrAuthTestNotImplemented
	}

	profile := request.Profile
	if profile == "" {
		activeProfile, err := service.profileStore.GetActiveProfile(ctx)
		if err != nil {
			return fmt.Errorf("read active profile: %w", err)
		}
		if activeProfile == "" {
			return ErrNoActiveProfile
		}
		profile = activeProfile
	}
	if err := domainauth.ValidateProfileName(profile); err != nil {
		return err
	}

	credential, err := service.credentialStore.Load(ctx, profile)
	if err != nil {
		return fmt.Errorf("load credential: %w", err)
	}
	defer domainauth.WipeBytes(credential.APISecret)

	if err := service.probe.Probe(ctx, credential); err != nil {
		return fmt.Errorf("probe auth: %w", err)
	}

	return nil
}
