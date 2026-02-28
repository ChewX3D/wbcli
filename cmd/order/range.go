package ordercmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

type rangeOptions struct {
	baseOptions
	StartPrice      float64
	EndPrice        float64
	Step            float64
	AmountMode      string
	BaseAmount      float64
	StartMultiplier float64
	StepMultiplier  float64
	Ratio           float64
	MaxMultiplier   float64
	DryRun          bool
	Confirm         bool
}

func newRangeCmd() *cobra.Command {
	options := &rangeOptions{}

	command := &cobra.Command{
		Use:   "range",
		Short: "Build or submit a range order plan",
		RunE: func(command *cobra.Command, args []string) error {
			if err := validateBase(options.baseOptions); err != nil {
				return err
			}
			if options.StartPrice <= 0 {
				return errors.New("--start-price must be greater than 0")
			}
			if options.EndPrice <= 0 {
				return errors.New("--end-price must be greater than 0")
			}
			if options.Step <= 0 {
				return errors.New("--step must be greater than 0")
			}
			if options.BaseAmount <= 0 {
				return errors.New("--base-amount must be greater than 0")
			}
			if err := validateAmountMode(options.AmountMode); err != nil {
				return err
			}

			_, err := fmt.Fprintf(
				command.OutOrStdout(),
				"wbcli order range is not implemented yet (profile=%s market=%s side=%s start=%g end=%g step=%g amount-mode=%s base-amount=%g dry-run=%t confirm=%t)\n",
				options.Profile,
				options.Market,
				options.Side,
				options.StartPrice,
				options.EndPrice,
				options.Step,
				options.AmountMode,
				options.BaseAmount,
				options.DryRun,
				options.Confirm,
			)
			return err
		},
	}

	addBaseFlags(command, &options.baseOptions)
	command.Flags().Float64Var(&options.StartPrice, "start-price", 0, "range start price")
	command.Flags().Float64Var(&options.EndPrice, "end-price", 0, "range end price")
	command.Flags().Float64Var(&options.Step, "step", 0, "price step")
	command.Flags().StringVar(
		&options.AmountMode,
		"amount-mode",
		"constant",
		"Amount sizing strategy per generated order (i is zero-based): constant=base-amount; arithmetic=base-amount*(start-multiplier+i*step-multiplier); geometric=base-amount*ratio^i; capped-geometric=min(base-amount*ratio^i, base-amount*max-multiplier); fibonacci=base-amount*fib(i+1); custom-list=explicit multipliers list.",
	)
	command.Flags().Float64Var(&options.BaseAmount, "base-amount", 0, "base amount per order")
	command.Flags().Float64Var(&options.StartMultiplier, "start-multiplier", 1, "arithmetic start multiplier")
	command.Flags().Float64Var(&options.StepMultiplier, "step-multiplier", 1, "arithmetic step multiplier")
	command.Flags().Float64Var(&options.Ratio, "ratio", 2, "geometric ratio")
	command.Flags().Float64Var(&options.MaxMultiplier, "max-multiplier", 0, "cap for capped-geometric mode")
	command.Flags().BoolVar(&options.DryRun, "dry-run", false, "preview plan without submitting orders")
	command.Flags().BoolVar(&options.Confirm, "confirm", false, "confirm live batch placement")

	return command
}

func validateAmountMode(mode string) error {
	switch mode {
	case "constant", "arithmetic", "geometric", "capped-geometric", "fibonacci", "custom-list":
		return nil
	default:
		return errors.New("--amount-mode must be one of: constant, arithmetic, geometric, capped-geometric, fibonacci, custom-list")
	}
}
