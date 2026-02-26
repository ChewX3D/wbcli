package cmd

import (
	"fmt"

	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
	"github.com/spf13/cobra"
)

type authTestOptions struct {
	Profile string
}

func newAuthTestCmd() *cobra.Command {
	options := &authTestOptions{}

	command := &cobra.Command{
		Use:   "test",
		Short: "Validate credentials for a profile",
		Example: "wbcli auth test --profile prod\n" +
			"wbcli auth test",
		RunE: func(command *cobra.Command, args []string) error {
			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			err = services.test.Execute(command.Context(), authservice.TestRequest{Profile: options.Profile})
			if err != nil {
				return mapAuthError(err)
			}

			_, err = fmt.Fprintln(command.OutOrStdout(), "auth test passed")
			return err
		},
	}

	command.Flags().StringVar(&options.Profile, "profile", "", "credential profile name (optional; uses active profile if omitted)")

	return command
}
