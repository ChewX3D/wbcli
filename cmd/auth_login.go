package cmd

import (
	"fmt"
	"os"

	authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"
	clitools "github.com/ChewX3D/wbcli/internal/cli"
	"github.com/spf13/cobra"
)

func newAuthLoginCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "login",
		Short: "Store credentials from stdin",
		Long:  "Store API key and API secret in secure OS keychain backend using stdin-only input.",
		Example: `  # Option 1: local inline values
  WBCLI_API_KEY='1' WBCLI_API_SECRET='2' sh -c 'printf "%s\n%s\n" "$WBCLI_API_KEY" "$WBCLI_API_SECRET"' | wbcli auth login

  # Option 2: CI job with env vars already injected
  sh -c 'printf "%s\n%s\n" "$WBCLI_API_KEY" "$WBCLI_API_SECRET"' | wbcli auth login

  # Option 3: local file with exactly two lines
  # line 1 = api key, line 2 = api secret
  cat /tmp/wbcli-auth.txt | wbcli auth login`,
		RunE: func(command *cobra.Command, args []string) error {
			if inputFile, ok := command.InOrStdin().(*os.File); ok && clitools.IsTerminalInput(inputFile) {
				return mapAuthError(clitools.ErrCredentialInputMissing)
			}

			credentials, err := clitools.ReadCredentialPairFromReader(command.InOrStdin(), 16*1024)
			if err != nil {
				return mapAuthError(err)
			}

			return runWithAuthServices(command, func(services *authServices) error {
				result, err := services.login.Execute(command.Context(), authservice.LoginRequest{
					APIKey:    credentials.APIKey,
					APISecret: credentials.APISecret,
				})
				if err != nil {
					return err
				}

				_, err = fmt.Fprintf(
					command.OutOrStdout(),
					"logged_in=true backend=%s api_key=%s saved_at=%s\n",
					result.Backend,
					result.APIKeyHint,
					result.SavedAt,
				)
				return err
			})
		},
	}

	return command
}
