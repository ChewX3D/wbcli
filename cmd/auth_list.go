package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newAuthListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List configured credential profiles",
		Example: "wbcli auth list",
		RunE: func(command *cobra.Command, args []string) error {
			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			result, err := services.list.Execute(command.Context())
			if err != nil {
				return mapAuthError(err)
			}
			if len(result.Profiles) == 0 {
				_, err := fmt.Fprintln(command.OutOrStdout(), "no profiles configured")
				return err
			}

			rows := []string{"PROFILE\tACTIVE\tBACKEND\tAPI_KEY\tUPDATED_AT"}
			for _, profile := range result.Profiles {
				active := ""
				if profile.Active {
					active = "*"
				}
				updated := "-"
				if !profile.UpdatedAt.IsZero() {
					updated = profile.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
				}
				rows = append(rows, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", profile.Name, active, profile.Backend, profile.APIKeyHint, updated))
			}

			_, err = fmt.Fprintln(command.OutOrStdout(), strings.Join(rows, "\n"))
			return err
		},
	}
}
