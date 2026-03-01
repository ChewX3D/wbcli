package ordercmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type placeOptions struct {
	baseOptions
	Amount        float64
	Price         float64
	Expiration    int64
	ClientOrderID string
}

func newPlaceCmd() *cobra.Command {
	options := &placeOptions{}

	command := &cobra.Command{
		Use:   "place",
		Short: "Place a single collateral limit order",
		RunE: func(command *cobra.Command, args []string) error {
			if err := validateBase(options.baseOptions); err != nil {
				return err
			}
			if err := validatePositiveFloatFlag("--amount", options.Amount); err != nil {
				return err
			}
			if err := validatePositiveFloatFlag("--price", options.Price); err != nil {
				return err
			}
			if err := validateNonNegativeInt64Flag("--expiration", options.Expiration); err != nil {
				return err
			}

			_, err := fmt.Fprintf(
				command.OutOrStdout(),
				"wbcli order place is not implemented yet (profile=%s market=%s side=%s amount=%g price=%g expiration=%d client-order-id=%s)\n",
				options.Profile,
				options.Market,
				options.Side,
				options.Amount,
				options.Price,
				options.Expiration,
				options.ClientOrderID,
			)
			return err
		},
	}

	addBaseFlags(command, &options.baseOptions)
	command.Flags().Float64Var(&options.Amount, "amount", 0, "order amount")
	command.Flags().Float64Var(&options.Price, "price", 0, "limit price")
	command.Flags().Int64Var(&options.Expiration, "expiration", 0, "order expiration")
	command.Flags().StringVar(&options.ClientOrderID, "client-order-id", "", "client order id")

	return command
}
