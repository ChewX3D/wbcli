package cmd

import (
	"fmt"

	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
	"github.com/spf13/cobra"
)

type authUseOptions struct {
	Profile string
}

func newAuthUseCmd() *cobra.Command {
	options := &authUseOptions{}

	command := &cobra.Command{
		Use:   "use",
		Short: "Select active profile",
		Example: "wbcli auth use --profile prod\n" +
			"wbcli auth use --profile sandbox",
		RunE: func(command *cobra.Command, args []string) error {
			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			result, err := services.use.Execute(command.Context(), authservice.UseRequest{Profile: options.Profile})
			if err != nil {
				return mapAuthError(err)
			}

			_, err = fmt.Fprintf(command.OutOrStdout(), "active profile set to %s\n", result.Profile)
			return err
		},
	}

	command.Flags().StringVar(&options.Profile, "profile", "", "credential profile name")
	_ = command.MarkFlagRequired("profile")

	return command
}
