package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show wbcli version",
		RunE: func(command *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(command.OutOrStdout(), "wbcli version %s\n", versionFromBuildConfig(getBuildConfig()))
			return err
		},
	}
}

func versionFromBuildConfig(raw []byte) string {
	const defaultVersion = "dev"

	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "version:") {
			continue
		}

		value := strings.TrimSpace(strings.TrimPrefix(trimmed, "version:"))
		if value == "" {
			return defaultVersion
		}

		return value
	}

	return defaultVersion
}
