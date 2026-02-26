package cmd

import (
	"fmt"
	"os"

	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
	clitools "github.com/ChewX3D/wbcli/internal/cli"
	"github.com/spf13/cobra"
)

type authLoginOptions struct {
	Profile string
	Force   bool
}

func newAuthLoginCmd() *cobra.Command {
	options := &authLoginOptions{}

	command := &cobra.Command{
		Use:   "login",
		Short: "Store credentials for a profile from stdin",
		Long:  "Store API key and API secret in secure OS keychain backend using stdin-only input.",
		Example: "printf '%s\\n%s\\n' \"$WBCLI_API_KEY\" \"$WBCLI_API_SECRET\" | wbcli auth login --profile prod\n" +
			"printf '%s\\n%s\\n' \"$WBCLI_API_KEY\" \"$WBCLI_API_SECRET\" | wbcli auth login --profile ci --force",
		RunE: func(command *cobra.Command, args []string) error {
			if inputFile, ok := command.InOrStdin().(*os.File); ok && clitools.IsTerminalInput(inputFile) {
				return mapAuthError(clitools.ErrCredentialInputMissing)
			}

			credentials, err := clitools.ReadCredentialPairFromReader(command.InOrStdin(), 16*1024)
			if err != nil {
				return mapAuthError(err)
			}

			services, err := authServicesFactory()
			if err != nil {
				return mapAuthError(err)
			}

			result, err := services.login.Execute(command.Context(), authservice.LoginRequest{
				Profile:   options.Profile,
				APIKey:    credentials.APIKey,
				APISecret: credentials.APISecret,
				Force:     options.Force,
			})
			if err != nil {
				return mapAuthError(err)
			}

			_, err = fmt.Fprintf(
				command.OutOrStdout(),
				"profile=%s backend=%s api_key=%s saved_at=%s\n",
				result.Profile,
				result.Backend,
				result.APIKeyHint,
				result.SavedAt,
			)
			return err
		},
	}

	command.Flags().StringVar(&options.Profile, "profile", "", "credential profile name")
	command.Flags().BoolVar(&options.Force, "force", false, "overwrite existing credentials for profile")
	_ = command.MarkFlagRequired("profile")

	return command
}
