package application

import (
	"context"
	"fmt"

	"github.com/ChewX3D/wbcli/internal/adapters/configstore"
	"github.com/ChewX3D/wbcli/internal/adapters/secretstore"
	"github.com/ChewX3D/wbcli/internal/adapters/whitebit"
	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
)

// AuthUseCases defines auth operations exposed to command adapters.
type AuthUseCases interface {
	Login(ctx context.Context, request authservice.LoginRequest) (authservice.LoginResult, error)
	Logout(ctx context.Context) (authservice.LogoutResult, error)
	Status(ctx context.Context) (authservice.StatusResult, error)
}

// Application holds use-case interfaces used by CLI command adapters.
type Application struct {
	Auth AuthUseCases
}

type authUseCases struct {
	login  *authservice.LoginService
	logout *authservice.LogoutService
	status *authservice.StatusService
}

// New constructs application container from prepared use-case interfaces.
func New(auth AuthUseCases) *Application {
	return &Application{Auth: auth}
}

// NewWithAuthServices constructs application container from concrete auth services.
func NewWithAuthServices(
	login *authservice.LoginService,
	logout *authservice.LogoutService,
	status *authservice.StatusService,
) *Application {
	return New(&authUseCases{
		login:  login,
		logout: logout,
		status: status,
	})
}

// NewDefault wires adapters and services for production runtime.
func NewDefault() (*Application, error) {
	sessionStore, err := configstore.NewDefaultSessionStore()
	if err != nil {
		return nil, fmt.Errorf("init session store: %w", err)
	}

	credentialStore := secretstore.NewOSKeychainStore()
	credentialVerifier := whitebit.NewDefaultCredentialVerifierAdapter()
	clock := authservice.SystemClock{}

	return NewWithAuthServices(
		authservice.NewLoginService(credentialStore, sessionStore, clock, credentialVerifier),
		authservice.NewLogoutService(credentialStore, sessionStore),
		authservice.NewStatusService(sessionStore),
	), nil
}

func (useCases *authUseCases) Login(ctx context.Context, request authservice.LoginRequest) (authservice.LoginResult, error) {
	return useCases.login.Execute(ctx, request)
}

func (useCases *authUseCases) Logout(ctx context.Context) (authservice.LogoutResult, error) {
	return useCases.logout.Execute(ctx)
}

func (useCases *authUseCases) Status(ctx context.Context) (authservice.StatusResult, error) {
	return useCases.status.Execute(ctx)
}
