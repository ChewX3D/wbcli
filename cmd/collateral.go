package cmd

import (
	collateralcmd "github.com/ChewX3D/wbcli/cmd/collateral"
	"github.com/spf13/cobra"
)

func newCollateralCmd(provider applicationProvider) *cobra.Command {
	return collateralcmd.NewCommand(provider)
}
