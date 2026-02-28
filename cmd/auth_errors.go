package cmd

import (
	"errors"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
	clitools "github.com/ChewX3D/wbcli/internal/cli"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

func mapAuthError(err error) error {
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
	case errors.Is(err, authservice.ErrCredentialAlreadyExists):
		return errors.New("credentials already exist; rerun with --force to overwrite")
	case errors.Is(err, authservice.ErrNotLoggedIn):
		return errors.New("not logged in; run wbcli auth login first")
	case errors.Is(err, authservice.ErrAuthTestNotImplemented):
		return errors.New("auth test is not implemented yet; it will be enabled after WhiteBIT client integration")
	case errors.Is(err, ports.ErrCredentialNotFound):
		return errors.New("not logged in; run wbcli auth login first")
	case errors.Is(err, ports.ErrSecretStoreUnavailable):
		return errors.New("os-keychain backend is unavailable on this system")
	case errors.Is(err, ports.ErrSecretStorePermissionDenied):
		return errors.New("os-keychain access denied; unlock keychain/secret service and retry")
	default:
		return err
	}
}
