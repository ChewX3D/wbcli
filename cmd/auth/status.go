package authcmd

import (
	"fmt"

	appcontainer "github.com/ChewX3D/wbcli/internal/app/application"
	"github.com/spf13/cobra"
)

func newStatusCmd(getApplication func() (*appcontainer.Application, error)) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Show current auth status",
		Example: "wbcli auth status",
		RunE: func(command *cobra.Command, args []string) error {
			return runWithApplication(command, getApplication, func(application *appcontainer.Application) error {
				result, err := application.Auth.Status(command.Context())
				if err != nil {
					return err
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
			})
		},
	}
}
