package cmd

import "github.com/spf13/cobra"

func newOrderCmd() *cobra.Command {
	orderCmd := &cobra.Command{
		Use:   "order",
		Short: "Place and manage orders",
		Long:  "Place single collateral orders or build/submit range order plans.",
		RunE: func(command *cobra.Command, args []string) error {
			return command.Help()
		},
	}

	orderCmd.AddCommand(newOrderPlaceCmd())
	orderCmd.AddCommand(newOrderRangeCmd())

	return orderCmd
}

type orderBaseOptions struct {
	Profile string
	Market  string
	Side    string
}

func addOrderBaseFlags(command *cobra.Command, options *orderBaseOptions) {
	command.Flags().StringVar(&options.Profile, "profile", "default", "credential profile name")
	command.Flags().StringVar(&options.Market, "market", "", "whitebit market pair (for example BTC_PERP)")
	command.Flags().StringVar(&options.Side, "side", "", "order side: buy or sell")
}
