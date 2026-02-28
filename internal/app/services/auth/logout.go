package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChewX3D/wbcli/internal/app/ports"
)

// LogoutResult contains safe logout result.
type LogoutResult struct {
	LoggedOut bool
}

// LogoutService deletes single-session credentials.
type LogoutService struct {
	credentialStore ports.CredentialStore
	sessionStore    ports.SessionStore
}

// NewLogoutService constructs LogoutService.
func NewLogoutService(credentialStore ports.CredentialStore, sessionStore ports.SessionStore) *LogoutService {
	return &LogoutService{
		credentialStore: credentialStore,
		sessionStore:    sessionStore,
	}
}

// Execute removes credentials and session metadata.
func (service *LogoutService) Execute(ctx context.Context) (LogoutResult, error) {
	if err := service.credentialStore.Delete(ctx); err != nil {
		if !errors.Is(err, ports.ErrCredentialNotFound) {
			return LogoutResult{}, fmt.Errorf("delete credential: %w", err)
		}
	}

	if err := service.sessionStore.ClearSession(ctx); err != nil {
		return LogoutResult{}, fmt.Errorf("clear session metadata: %w", err)
	}

	return LogoutResult{LoggedOut: true}, nil
}
