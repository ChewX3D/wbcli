package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "logout",
		Short:   "Remove stored credentials",
		Example: "wbcli auth logout",
		RunE: func(command *cobra.Command, args []string) error {
			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			_, err = services.logout.Execute(command.Context())
			if err != nil {
				return mapAuthError(err)
			}

			_, err = fmt.Fprintln(command.OutOrStdout(), "logged_out=true")
			return err
		},
	}
}
