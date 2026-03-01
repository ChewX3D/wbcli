package authcmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
	clitools "github.com/ChewX3D/wbcli/internal/cli"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domainauth.ErrAPIKeyRequired):
		return errors.New("api key is required in stdin payload")
	case errors.Is(err, domainauth.ErrAPISecretRequired):
		return errors.New("api secret is required in stdin payload")
	case errors.Is(err, clitools.ErrCredentialInputMissing):
		return errors.New("stdin credentials are required: first line API key, second line API secret")
	case errors.Is(err, clitools.ErrCredentialInputTooLarge):
		return errors.New("stdin credential payload exceeds maximum allowed size")
	case errors.Is(err, clitools.ErrCredentialInputFormat):
		return errors.New("stdin credential payload must contain exactly two non-empty lines: api_key then api_secret")
	case errors.Is(err, authservice.ErrNotLoggedIn):
		return errors.New("not logged in; run wbcli auth login first")
	case errors.Is(err, ports.ErrCredentialNotFound):
		return errors.New("not logged in; run wbcli auth login first")
	case errors.Is(err, ports.ErrSecretStoreUnavailable):
		return errors.New("os-keychain backend is unavailable on this system; install/unlock keychain backend and retry")
	case errors.Is(err, ports.ErrSecretStorePermissionDenied):
		return errors.New("os-keychain access denied; keychain is locked or access is restricted")
	case errors.Is(err, ports.ErrCredentialVerifyUnauthorized):
		return formatCredentialVerificationError(
			"whitebit credential verification failed: credentials are invalid",
			err,
		)
	case errors.Is(err, ports.ErrCredentialVerifyForbidden):
		return formatCredentialVerificationError(
			"whitebit credential verification failed: credentials are valid, but token permissions are insufficient for endpoint /api/v4/collateral-account/hedge-mode",
			err,
		)
	case errors.Is(err, ports.ErrCredentialVerifyUnavailable):
		return formatCredentialVerificationError(
			"whitebit credential verification unavailable: network issue or WhiteBIT service error; retry later",
			err,
		)
	default:
		return err
	}
}

func formatCredentialVerificationError(baseMessage string, err error) error {
	reason := extractCredentialVerificationReason(err)
	if reason == "" {
		return errors.New(baseMessage)
	}

	return fmt.Errorf("%s. reason: %s", baseMessage, reason)
}

func extractCredentialVerificationReason(err error) string {
	if err == nil {
		return ""
	}

	reason := err.Error()
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return ""
	}

	prefixes := []string{
		"verify credential: ",
		"credential verification unauthorized: ",
		"credential verification forbidden: ",
		"credential verification unavailable: ",
	}

	for {
		updated := reason
		for _, prefix := range prefixes {
			updated = strings.TrimPrefix(updated, prefix)
		}
		if updated == reason {
			break
		}
		reason = updated
	}

	return strings.TrimSpace(reason)
}
