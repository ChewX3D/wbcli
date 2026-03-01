package cmd

import (
	authcmd "github.com/ChewX3D/wbcli/cmd/auth"
	"github.com/spf13/cobra"
)

func newAuthCmd(provider applicationProvider) *cobra.Command {
	return authcmd.NewCommand(provider)
}
