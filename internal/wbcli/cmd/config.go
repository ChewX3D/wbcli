package cmd

var buildConfig []byte

func SetBuildConfig(cfg []byte) {
	buildConfig = make([]byte, len(cfg))
	copy(buildConfig, cfg)
}

func getBuildConfig() []byte {
	cfg := make([]byte, len(buildConfig))
	copy(cfg, buildConfig)

	return cfg
}
