package cmd

import "github.com/spf13/cobra"

func runWithAuthServices(command *cobra.Command, run func(*authServices) error) error {
	services, err := authServicesFactory()
	if err != nil {
		return mapAuthError(err)
	}

	if err := run(services); err != nil {
		return mapAuthError(err)
	}

	return nil
}
