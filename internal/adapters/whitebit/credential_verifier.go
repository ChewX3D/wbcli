package whitebit

import (
	"context"
	"errors"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

// CredentialVerifierAdapter adapts auth credential verification port to WhiteBIT client endpoints.
type CredentialVerifierAdapter struct {
	client *Client
}

// NewCredentialVerifierAdapter constructs CredentialVerifierAdapter.
func NewCredentialVerifierAdapter(client *Client) *CredentialVerifierAdapter {
	return &CredentialVerifierAdapter{client: client}
}

// NewDefaultCredentialVerifierAdapter constructs CredentialVerifierAdapter with default WhiteBIT client.
func NewDefaultCredentialVerifierAdapter() *CredentialVerifierAdapter {
	return NewCredentialVerifierAdapter(NewDefaultClient())
}

// Verify checks login credentials using the documented hedge-mode endpoint.
func (adapter *CredentialVerifierAdapter) Verify(ctx context.Context, credential domainauth.Credential) error {
	if adapter == nil || adapter.client == nil {
		return ports.ErrCredentialVerifyUnavailable
	}

	_, err := adapter.client.GetCollateralAccountHedgeMode(ctx, credential)
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrUnauthorized):
		return ports.ErrCredentialVerifyUnauthorized
	case errors.Is(err, ErrForbidden):
		return ports.ErrCredentialVerifyForbidden
	default:
		return ports.ErrCredentialVerifyUnavailable
	}
}
