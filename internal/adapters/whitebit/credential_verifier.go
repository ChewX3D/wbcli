package whitebit

import (
	"context"

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
func (adapter *CredentialVerifierAdapter) Verify(ctx context.Context, credential domainauth.Credential) (ports.CredentialVerificationResult, error) {
	if adapter == nil || adapter.client == nil {
		return ports.CredentialVerificationResult{}, &ports.APIError{
			Code:    ports.CodeUnavailable,
			Message: "credential verification failed: exchange unavailable",
			Details: "credential verifier adapter is not configured",
		}
	}

	response, err := adapter.client.GetCollateralAccountHedgeMode(ctx, credential)
	if err == nil {
		return ports.CredentialVerificationResult{
			Endpoint:  collateralAccountHedgeModePath,
			HedgeMode: response.HedgeMode,
		}, nil
	}

	return ports.CredentialVerificationResult{}, buildAPIError(err, collateralAccountHedgeModePath, "credential verification")
}
