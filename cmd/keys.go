package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:     "auth",
		Aliases: []string{"keys"},
		Short:   "Manage authentication credentials",
		Long:    "Manage WhiteBIT API authentication credentials by profile.",
		RunE: func(command *cobra.Command, args []string) error {
			return command.Help()
		},
	}

	authCmd.AddCommand(newKeysSetCmd())
	authCmd.AddCommand(newKeysListCmd())
	authCmd.AddCommand(newKeysRemoveCmd())
	authCmd.AddCommand(newKeysTestCmd())

	return authCmd
}

type authProfileOptions struct {
	Profile string
}

func addProfileFlag(command *cobra.Command, options *authProfileOptions) {
	command.Flags().StringVar(&options.Profile, "profile", "default", "credential profile name")
}

func writeNotImplemented(command *cobra.Command, action string, profile string) error {
	_, err := fmt.Fprintf(command.OutOrStdout(), "wbcli auth %s is not implemented yet (profile=%s)\n", action, profile)
	return err
}
