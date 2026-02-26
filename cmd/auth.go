package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication credentials",
		Long:  "Manage WhiteBIT API authentication credentials by profile.",
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("unknown command %q for %q", args[0], command.CommandPath())
			}

			return command.Help()
		},
	}

	authCmd.AddCommand(newAuthLoginCmd())
	authCmd.AddCommand(newAuthUseCmd())
	authCmd.AddCommand(newAuthListCmd())
	authCmd.AddCommand(newAuthLogoutCmd())
	authCmd.AddCommand(newAuthCurrentCmd())
	authCmd.AddCommand(newAuthTestCmd())

	return authCmd
}
