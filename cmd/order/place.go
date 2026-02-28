package ordercmd

import (
	"errors"
	"fmt"
	"strings"

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
			if options.Amount <= 0 {
				return errors.New("--amount must be greater than 0")
			}
			if options.Price <= 0 {
				return errors.New("--price must be greater than 0")
			}
			if options.Expiration < 0 {
				return errors.New("--expiration must be greater than or equal to 0")
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

func validateBase(options baseOptions) error {
	if options.Market == "" {
		return errors.New("--market is required")
	}

	side := strings.ToLower(options.Side)
	if side == "" {
		return errors.New("--side is required")
	}
	if side != "buy" && side != "sell" {
		return errors.New("--side must be one of: buy, sell")
	}

	return nil
}
