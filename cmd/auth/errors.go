package authcmd

import (
	"errors"

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
		return errors.New("whitebit credential verification failed: credentials are invalid (wrong public/secret key pair, disabled key, or signature/nonce rejected)")
	case errors.Is(err, ports.ErrCredentialVerifyForbidden):
		return errors.New("whitebit credential verification failed: credentials are valid, but token permissions are insufficient for endpoint /api/v4/collateral-account/hedge-mode")
	case errors.Is(err, ports.ErrCredentialVerifyUnavailable):
		return errors.New("whitebit credential verification unavailable: network issue or WhiteBIT service error; retry later")
	default:
		return err
	}
}
