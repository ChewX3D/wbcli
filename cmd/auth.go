package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication credentials",
		Long:  "Manage WhiteBIT API authentication credentials in single-session mode.",
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("unknown command %q for %q", args[0], command.CommandPath())
			}

			return command.Help()
		},
	}

	authCmd.AddCommand(newAuthLoginCmd())
	authCmd.AddCommand(newAuthLogoutCmd())
	authCmd.AddCommand(newAuthStatusCmd())
	authCmd.AddCommand(newAuthTestCmd())

	return authCmd
}
