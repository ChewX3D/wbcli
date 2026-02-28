package authcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "logout",
		Short:   "Remove stored credentials",
		Example: "wbcli auth logout",
		RunE: func(command *cobra.Command, args []string) error {
			return runWithServices(command, func(services *Services) error {
				if _, err := services.Logout.Execute(command.Context()); err != nil {
					return err
				}

				_, err := fmt.Fprintln(command.OutOrStdout(), "logged_out=true")
				return err
			})
		},
	}
}
