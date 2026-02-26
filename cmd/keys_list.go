package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newKeysListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured credential profiles",
		RunE: func(command *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(command.OutOrStdout(), "wbcli auth list is not implemented yet")
			return err
		},
	}
}
