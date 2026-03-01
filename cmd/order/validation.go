package ordercmd

import (
	"errors"
	"fmt"
	"strings"
)

const amountModeErrorMessage = "--amount-mode must be one of: constant, arithmetic, geometric, capped-geometric, fibonacci, custom-list"

var supportedAmountModes = map[string]struct{}{
	"constant":         {},
	"arithmetic":       {},
	"geometric":        {},
	"capped-geometric": {},
	"fibonacci":        {},
	"custom-list":      {},
}

func validatePositiveFloatFlag(flagName string, value float64) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than 0", flagName)
	}

	return nil
}

func validateNonNegativeInt64Flag(flagName string, value int64) error {
	if value < 0 {
		return fmt.Errorf("%s must be greater than or equal to 0", flagName)
	}

	return nil
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

func validateAmountMode(mode string) error {
	if _, ok := supportedAmountModes[mode]; ok {
		return nil
	}

	return errors.New(amountModeErrorMessage)
}
