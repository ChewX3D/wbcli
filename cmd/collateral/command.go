package collateralcmd

import (
	ordercmd "github.com/ChewX3D/wbcli/cmd/order"
	appcontainer "github.com/ChewX3D/wbcli/internal/app/application"
	"github.com/spf13/cobra"
)

// NewCommand constructs collateral command group.
func NewCommand(provider func() (*appcontainer.Application, error)) *cobra.Command {
	command := &cobra.Command{
		Use:   "collateral",
		Short: "Collateral trading commands",
		Long:  "Run collateral trading workflows such as single order placement and range planning.",
		RunE: func(command *cobra.Command, args []string) error {
			return command.Help()
		},
	}

	command.AddCommand(ordercmd.NewCommand(provider))

	return command
}
