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
		Example: `  # Local shell with env vars
  printf '%s\n%s\n' "$WBCLI_API_KEY" "$WBCLI_API_SECRET" | wbcli auth login

  # CI job with secrets injected as env vars
  printf '%s\n%s\n' "$WBCLI_API_KEY" "$WBCLI_API_SECRET" | wbcli auth login

  # Read from a local file that contains two lines:
  # line 1 = api key, line 2 = api secret
  cat ./secrets/wbcli-auth.txt | wbcli auth login`,
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
