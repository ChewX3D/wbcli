package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newKeysCmd() *cobra.Command {
	keysCmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage API credential profiles",
		Long:  "Manage API credentials by profile for secure WhiteBIT API usage.",
		RunE: func(command *cobra.Command, args []string) error {
			return command.Help()
		},
	}

	keysCmd.AddCommand(newKeysSetCmd())
	keysCmd.AddCommand(newKeysListCmd())
	keysCmd.AddCommand(newKeysRemoveCmd())
	keysCmd.AddCommand(newKeysTestCmd())

	return keysCmd
}

type keysProfileOptions struct {
	Profile string
}

func addProfileFlag(command *cobra.Command, options *keysProfileOptions) {
	command.Flags().StringVar(&options.Profile, "profile", "default", "credential profile name")
}

func writeNotImplemented(command *cobra.Command, action string, profile string) error {
	_, err := fmt.Fprintf(command.OutOrStdout(), "wbcli keys %s is not implemented yet (profile=%s)\n", action, profile)
	return err
}
