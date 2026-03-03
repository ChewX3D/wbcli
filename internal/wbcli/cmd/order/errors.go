package ordercmd

import (
	"errors"

	"github.com/ChewX3D/crypto/internal/app/ports"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr *ports.APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	switch {
	case errors.Is(err, ports.ErrCredentialNotFound):
		return errors.New("not logged in; run wbcli auth login first")
	case errors.Is(err, ports.ErrSecretStoreUnavailable):
		return errors.New("os-keychain backend is unavailable on this system; install/unlock keychain backend and retry")
	case errors.Is(err, ports.ErrSecretStorePermissionDenied):
		return errors.New("os-keychain access denied; keychain is locked or access is restricted")
	}

	return err
}
