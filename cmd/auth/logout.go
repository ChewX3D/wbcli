package authcmd

import (
	"fmt"

	appcontainer "github.com/ChewX3D/wbcli/internal/app/application"
	"github.com/spf13/cobra"
)

func newLogoutCmd(getApplication func() (*appcontainer.Application, error)) *cobra.Command {
	return &cobra.Command{
		Use:     "logout",
		Short:   "Remove stored credentials",
		Example: "wbcli auth logout",
		RunE: func(command *cobra.Command, args []string) error {
			return runWithApplication(command, getApplication, func(application *appcontainer.Application) error {
				if _, err := application.Auth.Logout(command.Context()); err != nil {
					return err
				}

				_, err := fmt.Fprintln(command.OutOrStdout(), "logged_out=true")
				return err
			})
		},
	}
}
