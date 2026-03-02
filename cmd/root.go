package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

const flagKeyVerbose = "verbose"

var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	applicationProvider := newApplicationProvider()

	root := &cobra.Command{
		Use:   "wbcli",
		Short: "A safe CLI for WhiteBIT trading workflows",
		Long: `wbcli is a CLI for WhiteBIT collateral trading workflows.
It provides safe command groups for auth credential management and collateral order execution.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			isVerbose, err := cmd.Flags().GetBool(flagKeyVerbose)
			if err != nil {
				return err
			}

			if isVerbose {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}

			return nil
		},
	}

	root.PersistentFlags().BoolP(flagKeyVerbose, "v", false, "verbose logging")
	root.AddCommand(newVersionCmd())
	root.AddCommand(newAuthCmd(applicationProvider))
	root.AddCommand(newCollateralCmd(applicationProvider))

	return root
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func NewRootCmdForTest() *cobra.Command {
	return newRootCmd()
}
