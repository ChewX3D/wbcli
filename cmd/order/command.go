package ordercmd

import (
	appcontainer "github.com/ChewX3D/wbcli/internal/app/application"
	"github.com/spf13/cobra"
)

// NewCommand constructs the order command group.
func NewCommand(_ func() (*appcontainer.Application, error)) *cobra.Command {
	orderCmd := &cobra.Command{
		Use:   "order",
		Short: "Place and manage orders",
		Long:  "Place single collateral orders or build/submit range order plans.",
		RunE: func(command *cobra.Command, args []string) error {
			return command.Help()
		},
	}

	orderCmd.AddCommand(newPlaceCmd())
	orderCmd.AddCommand(newRangeCmd())

	return orderCmd
}

type baseOptions struct {
	Profile string
	Market  string
	Side    string
}

func addBaseFlags(command *cobra.Command, options *baseOptions) {
	command.Flags().StringVar(&options.Profile, "profile", "default", "credential profile name")
	command.Flags().StringVar(&options.Market, "market", "", "whitebit market pair (for example BTC_PERP)")
	command.Flags().StringVar(&options.Side, "side", "", "order side: buy or sell")
}
