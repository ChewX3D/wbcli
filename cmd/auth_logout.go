package cmd

import (
	"fmt"

	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
	"github.com/spf13/cobra"
)

type authLogoutOptions struct {
	Profile string
}

func newAuthLogoutCmd() *cobra.Command {
	options := &authLogoutOptions{}

	command := &cobra.Command{
		Use:   "logout",
		Short: "Remove credentials for a profile",
		Example: "wbcli auth logout --profile prod\n" +
			"wbcli auth logout --profile sandbox",
		RunE: func(command *cobra.Command, args []string) error {
			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			result, err := services.logout.Execute(command.Context(), authservice.LogoutRequest{Profile: options.Profile})
			if err != nil {
				return mapAuthError(err)
			}

			_, err = fmt.Fprintf(command.OutOrStdout(), "profile %s logged out\n", result.Profile)
			return err
		},
	}

	command.Flags().StringVar(&options.Profile, "profile", "", "credential profile name")
	_ = command.MarkFlagRequired("profile")

	return command
}
