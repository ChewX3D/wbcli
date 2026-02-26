package cmd

import "github.com/spf13/cobra"

func newKeysTestCmd() *cobra.Command {
	options := &authProfileOptions{}

	command := &cobra.Command{
		Use:   "test",
		Short: "Validate credentials for a profile",
		RunE: func(command *cobra.Command, args []string) error {
			return writeNotImplemented(command, "test", options.Profile)
		},
	}

	addProfileFlag(command, options)

	return command
}
