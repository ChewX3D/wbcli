package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

// TestService validates credentials against WhiteBIT via probe port.
type TestService struct {
	credentialStore ports.CredentialStore
	probe           ports.AuthProbe
}

// NewTestService constructs TestService.
func NewTestService(credentialStore ports.CredentialStore, probe ports.AuthProbe) *TestService {
	return &TestService{
		credentialStore: credentialStore,
		probe:           probe,
	}
}

// Execute runs auth test when probe is available.
func (service *TestService) Execute(ctx context.Context) error {
	if service.probe == nil {
		return ErrAuthTestNotImplemented
	}

	credential, err := service.credentialStore.Load(ctx)
	if err != nil {
		if errors.Is(err, ports.ErrCredentialNotFound) {
			return ErrNotLoggedIn
		}
		return fmt.Errorf("load credential: %w", err)
	}
	defer domainauth.WipeBytes(credential.APISecret)

	if err := service.probe.Probe(ctx, credential); err != nil {
		return fmt.Errorf("probe auth: %w", err)
	}

	return nil
}
