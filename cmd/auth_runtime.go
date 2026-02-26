package cmd

import (
	"fmt"

	"github.com/ChewX3D/wbcli/internal/adapters/configstore"
	"github.com/ChewX3D/wbcli/internal/adapters/secretstore"
	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
)

type authServices struct {
	login   *authservice.LoginService
	use     *authservice.UseService
	list    *authservice.ListService
	logout  *authservice.LogoutService
	current *authservice.CurrentService
	test    *authservice.TestService
}

var authServicesFactory = defaultAuthServicesFactory

func defaultAuthServicesFactory() (*authServices, error) {
	profileStore, err := configstore.NewDefaultProfileStore()
	if err != nil {
		return nil, fmt.Errorf("init profile store: %w", err)
	}

	credentialStore := secretstore.NewOSKeychainStore()
	clock := authservice.SystemClock{}

	return &authServices{
		login:   authservice.NewLoginService(credentialStore, profileStore, clock),
		use:     authservice.NewUseService(credentialStore, profileStore, clock),
		list:    authservice.NewListService(profileStore),
		logout:  authservice.NewLogoutService(credentialStore, profileStore),
		current: authservice.NewCurrentService(profileStore),
		test:    authservice.NewTestService(credentialStore, profileStore, nil),
	}, nil
}
