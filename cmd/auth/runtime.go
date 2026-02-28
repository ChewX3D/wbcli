package authcmd

import (
	"fmt"

	"github.com/ChewX3D/wbcli/internal/adapters/configstore"
	"github.com/ChewX3D/wbcli/internal/adapters/secretstore"
	"github.com/ChewX3D/wbcli/internal/adapters/whitebit"
	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
)

// Services contains auth command dependencies.
type Services struct {
	Login  *authservice.LoginService
	Logout *authservice.LogoutService
	Status *authservice.StatusService
}

var servicesFactory = defaultServicesFactory

func defaultServicesFactory() (*Services, error) {
	sessionStore, err := configstore.NewDefaultSessionStore()
	if err != nil {
		return nil, fmt.Errorf("init session store: %w", err)
	}

	credentialStore := secretstore.NewOSKeychainStore()
	credentialVerifier := whitebit.NewDefaultClient()
	clock := authservice.SystemClock{}

	return &Services{
		Login:  authservice.NewLoginService(credentialStore, sessionStore, clock, credentialVerifier),
		Logout: authservice.NewLogoutService(credentialStore, sessionStore),
		Status: authservice.NewStatusService(sessionStore),
	}, nil
}

// SetServicesFactoryForTest overrides auth service wiring in tests.
func SetServicesFactoryForTest(factory func() (*Services, error)) func() {
	previousFactory := servicesFactory
	servicesFactory = factory

	return func() {
		servicesFactory = previousFactory
	}
}
