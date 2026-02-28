package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
)

// StatusResult returns current auth session state.
type StatusResult struct {
	LoggedIn   bool
	Backend    string
	APIKeyHint string
	UpdatedAt  time.Time
}

// StatusService provides current auth status view.
type StatusService struct {
	sessionStore ports.SessionStore
}

// NewStatusService constructs StatusService.
func NewStatusService(sessionStore ports.SessionStore) *StatusService {
	return &StatusService{sessionStore: sessionStore}
}

// Execute returns logged-in/logged-out status with safe metadata.
func (service *StatusService) Execute(ctx context.Context) (StatusResult, error) {
	session, found, err := service.sessionStore.GetSession(ctx)
	if err != nil {
		return StatusResult{}, fmt.Errorf("read session metadata: %w", err)
	}
	if !found {
		return StatusResult{LoggedIn: false}, nil
	}

	return StatusResult{
		LoggedIn:   true,
		Backend:    session.Backend,
		APIKeyHint: session.APIKeyHint,
		UpdatedAt:  session.UpdatedAt,
	}, nil
}
