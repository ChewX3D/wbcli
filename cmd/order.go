package cmd

import (
	ordercmd "github.com/ChewX3D/wbcli/cmd/order"
	"github.com/spf13/cobra"
)

func newOrderCmd() *cobra.Command { return ordercmd.NewCommand() }
