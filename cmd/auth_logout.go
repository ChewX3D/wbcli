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
			return runWithAuthServices(command, func(services *authServices) error {
				if _, err := services.logout.Execute(command.Context()); err != nil {
					return err
				}

				_, err := fmt.Fprintln(command.OutOrStdout(), "logged_out=true")
				return err
			})
		},
	}
}
