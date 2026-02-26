package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "current",
		Short:   "Show current active profile",
		Example: "wbcli auth current",
		RunE: func(command *cobra.Command, args []string) error {
			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			result, err := services.current.Execute(command.Context())
			if err != nil {
				return mapAuthError(err)
			}

			_, err = fmt.Fprintf(
				command.OutOrStdout(),
				"profile=%s backend=%s api_key=%s updated_at=%s\n",
				result.Profile,
				result.Backend,
				result.APIKeyHint,
				result.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
			)
			return err
		},
	}
}
