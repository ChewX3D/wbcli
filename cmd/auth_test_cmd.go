package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "test",
		Short:   "Validate stored credentials",
		Example: "wbcli auth test",
		RunE: func(command *cobra.Command, args []string) error {
			return runWithAuthServices(command, func(services *authServices) error {
				if err := services.test.Execute(command.Context()); err != nil {
					return err
				}

				_, err := fmt.Fprintln(command.OutOrStdout(), "auth test passed")
				return err
			})
		},
	}
}
