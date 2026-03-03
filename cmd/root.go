package cmd

import (
	"log/slog"
	"os"

	appcontainer "github.com/ChewX3D/wbcli/internal/app/application"
	"github.com/spf13/cobra"
)

const flagKeyVerbose = "verbose"

func newRootCmd(factory func() (*appcontainer.Application, error)) *cobra.Command {
	applicationProvider := newApplicationProvider(factory)

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

// Execute creates the root command with production defaults and runs it.
func Execute() {
	if err := newRootCmd(appcontainer.NewDefault).Execute(); err != nil {
		os.Exit(1)
	}
}

// NewRootCmdForTest creates a root command with the given factory for tests.
func NewRootCmdForTest(factory func() (*appcontainer.Application, error)) *cobra.Command {
	return newRootCmd(factory)
}
