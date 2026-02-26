package main

import (
	_ "embed"

	"github.com/ChewX3D/wbcli/cmd"
)

//go:embed cmd/wbcli/config.yaml
var buildConfig []byte

func main() {
	cmd.SetBuildConfig(buildConfig)
	cmd.Execute()
}
