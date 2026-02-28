package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show current auth status",
		Example: "wbcli auth status",
		RunE: func(command *cobra.Command, args []string) error {
			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			result, err := services.status.Execute(command.Context())
			if err != nil {
				return mapAuthError(err)
			}
			if !result.LoggedIn {
				_, err := fmt.Fprintln(command.OutOrStdout(), "logged_in=false")
				return err
			}

			_, err = fmt.Fprintf(
				command.OutOrStdout(),
				"logged_in=true backend=%s api_key=%s updated_at=%s\n",
				result.Backend,
				result.APIKeyHint,
				result.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
			)
			return err
		},
	}
}
