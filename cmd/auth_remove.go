package cmd

import "github.com/spf13/cobra"

func newAuthRemoveCmd() *cobra.Command {
	options := &authProfileOptions{}

	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove credentials for a profile",
		RunE: func(command *cobra.Command, args []string) error {
			return writeNotImplemented(command, "remove", options.Profile)
		},
	}

	addProfileFlag(command, options)

	return command
}
