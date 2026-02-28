package cmd

import (
	"fmt"

	"github.com/ChewX3D/wbcli/internal/adapters/configstore"
	"github.com/ChewX3D/wbcli/internal/adapters/secretstore"
	"github.com/ChewX3D/wbcli/internal/adapters/whitebit"
	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
)

type authServices struct {
	login  *authservice.LoginService
	logout *authservice.LogoutService
	status *authservice.StatusService
}

var authServicesFactory = defaultAuthServicesFactory

func defaultAuthServicesFactory() (*authServices, error) {
	sessionStore, err := configstore.NewDefaultSessionStore()
	if err != nil {
		return nil, fmt.Errorf("init session store: %w", err)
	}

	credentialStore := secretstore.NewOSKeychainStore()
	authProbe := whitebit.NewDefaultAuthProbe()
	clock := authservice.SystemClock{}

	return &authServices{
		login:  authservice.NewLoginService(credentialStore, sessionStore, clock, authProbe),
		logout: authservice.NewLogoutService(credentialStore, sessionStore),
		status: authservice.NewStatusService(sessionStore),
	}, nil
}
