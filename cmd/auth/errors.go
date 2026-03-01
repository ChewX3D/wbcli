package authcmd

import (
	"errors"
	"fmt"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
	clitools "github.com/ChewX3D/wbcli/internal/cli"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

type staticAuthErrorRule struct {
	match   error
	message string
}

var staticAuthErrorRules = []staticAuthErrorRule{
	{match: domainauth.ErrAPIKeyRequired, message: "api key is required in stdin payload"},
	{match: domainauth.ErrAPISecretRequired, message: "api secret is required in stdin payload"},
	{match: clitools.ErrCredentialInputMissing, message: "stdin credentials are required: first line API key, second line API secret"},
	{match: clitools.ErrCredentialInputTooLarge, message: "stdin credential payload exceeds maximum allowed size"},
	{match: clitools.ErrCredentialInputFormat, message: "stdin credential payload must contain exactly two non-empty lines: api_key then api_secret"},
	{match: authservice.ErrNotLoggedIn, message: "not logged in; run wbcli auth login first"},
	{match: ports.ErrCredentialNotFound, message: "not logged in; run wbcli auth login first"},
	{match: ports.ErrSecretStoreUnavailable, message: "os-keychain backend is unavailable on this system; install/unlock keychain backend and retry"},
	{match: ports.ErrSecretStorePermissionDenied, message: "os-keychain access denied; keychain is locked or access is restricted"},
	{match: ports.ErrCredentialVerifyUnauthorized, message: "whitebit credential verification failed: credentials are invalid"},
	{match: ports.ErrCredentialVerifyForbidden, message: "whitebit credential verification failed: token permissions are insufficient"},
	{match: ports.ErrCredentialVerifyUnavailable, message: "whitebit credential verification unavailable: network issue or WhiteBIT service error; retry later"},
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	if mappedErr, ok := mapCredentialVerificationError(err); ok {
		return mappedErr
	}

	for _, rule := range staticAuthErrorRules {
		if errors.Is(err, rule.match) {
			return errors.New(rule.message)
		}
	}

	return err
}

func mapCredentialVerificationError(err error) (error, bool) {
	var verificationErr *ports.CredentialVerificationError
	if !errors.As(err, &verificationErr) {
		return nil, false
	}

	baseMessage := ""
	switch verificationErr.Reason {
	case ports.CredentialVerificationInvalidCredentials:
		baseMessage = "whitebit credential verification failed: credentials are invalid"
	case ports.CredentialVerificationInsufficientAccess:
		baseMessage = fmt.Sprintf(
			"whitebit credential verification failed: token permissions are insufficient; enable access to endpoint %s",
			normalizeEndpoint(verificationErr.Endpoint),
		)
	case ports.CredentialVerificationUnavailable:
		baseMessage = "whitebit credential verification unavailable: network issue or WhiteBIT service error; retry later"
	default:
		baseMessage = "whitebit credential verification failed"
	}

	if verificationErr.Detail == "" {
		return errors.New(baseMessage), true
	}

	return fmt.Errorf("%s. reason: %s", baseMessage, verificationErr.Detail), true
}

func normalizeEndpoint(endpoint string) string {
	if endpoint == "" {
		return "<required-endpoint>"
	}

	return endpoint
}
