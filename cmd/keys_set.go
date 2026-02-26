package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

type keysSetOptions struct {
	keysProfileOptions
	APIKey    string
	APISecret string
}

func newKeysSetCmd() *cobra.Command {
	options := &keysSetOptions{}

	command := &cobra.Command{
		Use:   "set",
		Short: "Store credentials for a profile",
		RunE: func(command *cobra.Command, args []string) error {
			if options.APIKey == "" {
				return errors.New("--api-key is required")
			}
			if options.APISecret == "" {
				return errors.New("--api-secret is required")
			}

			return writeNotImplemented(command, "set", options.Profile)
		},
	}

	addProfileFlag(command, &options.keysProfileOptions)
	command.Flags().StringVar(&options.APIKey, "api-key", "", "whitebit api key")
	command.Flags().StringVar(&options.APISecret, "api-secret", "", "whitebit api secret")

	return command
}
