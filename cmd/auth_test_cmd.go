package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "test",
		Short:   "Validate stored credentials",
		Example: "wbcli auth test",
		RunE: func(command *cobra.Command, args []string) error {
			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			err = services.test.Execute(command.Context())
			if err != nil {
				return mapAuthError(err)
			}

			_, err = fmt.Fprintln(command.OutOrStdout(), "auth test passed")
			return err
		},
	}
}
