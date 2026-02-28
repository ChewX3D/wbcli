package authcmd

import "github.com/spf13/cobra"

func runWithServices(command *cobra.Command, run func(*Services) error) error {
	services, err := servicesFactory()
	if err != nil {
		return mapError(err)
	}

	if err := run(services); err != nil {
		return mapError(err)
	}

	return nil
}
