package ordercmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ChewX3D/wbcli/internal/adapters/whitebit"
	"github.com/ChewX3D/wbcli/internal/app/ports"
)

const collateralLimitOrderEndpoint = "/api/v4/order/collateral/limit"

func mapError(err error) error {
	if err == nil {
		return nil
	}

	if mapped, ok := mapOrderAuthError(err); ok {
		return mapped
	}

	return err
}

func mapOrderAuthError(err error) (error, bool) {
	switch {
	case errors.Is(err, ports.ErrCredentialNotFound):
		return errors.New("not logged in; run wbcli auth login first"), true
	case errors.Is(err, ports.ErrSecretStoreUnavailable):
		return errors.New("os-keychain backend is unavailable on this system; install/unlock keychain backend and retry"), true
	case errors.Is(err, ports.ErrSecretStorePermissionDenied):
		return errors.New("os-keychain access denied; keychain is locked or access is restricted"), true
	case errors.Is(err, whitebit.ErrForbidden):
		return insufficientPermissionError(err), true
	case errors.Is(err, whitebit.ErrUnauthorized):
		if indicatesMissingEndpointAccess(err) {
			return insufficientPermissionError(err), true
		}

		baseMessage := "whitebit order placement failed: credentials are invalid"
		detail := extractErrorDetail(err)
		if detail == "" {
			return errors.New(baseMessage), true
		}

		return fmt.Errorf("%s. reason: %s", baseMessage, detail), true
	default:
		return nil, false
	}
}

func insufficientPermissionError(err error) error {
	baseMessage := fmt.Sprintf(
		"whitebit order placement failed: token permissions are insufficient; enable access to endpoint %s",
		collateralLimitOrderEndpoint,
	)
	detail := extractErrorDetail(err)
	if detail == "" {
		return errors.New(baseMessage)
	}

	return fmt.Errorf("%s. reason: %s", baseMessage, detail)
}

func extractErrorDetail(err error) string {
	if err == nil {
		return ""
	}

	detail := strings.TrimSpace(err.Error())
	if detail == "" {
		return ""
	}

	suffixes := []string{
		": whitebit unauthorized",
		": whitebit forbidden",
	}
	for _, suffix := range suffixes {
		detail = strings.TrimSuffix(detail, suffix)
	}

	return strings.TrimSpace(detail)
}

func indicatesMissingEndpointAccess(err error) bool {
	detail := strings.ToLower(extractErrorDetail(err))

	return strings.Contains(detail, "not authorized to perform this action") ||
		strings.Contains(detail, "not authorised to perform this action")
}
