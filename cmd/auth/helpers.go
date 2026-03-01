package authcmd

import (
	appcontainer "github.com/ChewX3D/wbcli/internal/app/application"
	"github.com/spf13/cobra"
)

func runWithApplication(
	command *cobra.Command,
	getApplication func() (*appcontainer.Application, error),
	run func(*appcontainer.Application) error,
) error {
	application, err := getApplication()
	if err != nil {
		return mapError(err)
	}

	if err := run(application); err != nil {
		return mapError(err)
	}

	return nil
}
