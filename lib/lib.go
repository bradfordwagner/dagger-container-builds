package lib

import (
	"context"
	"dagger/container-builds/internal/dagger"
)

// OpenConfigYaml opens the config.yaml file
func OpenConfigYaml(
	ctx context.Context,
	src dagger.Directory,
) (s string, err error) {
	configFile := src.File("config.yaml")
	return configFile.Contents(ctx)
}

func ArchToRunner(arch string) (s string) {
	archs := map[string]string{
		"linux/arm64": "arm64",
	}

	// default to ubuntu-latest
	s, ok := archs[arch]
	if !ok {
		return "ubuntu-latest"
	}
	return
}
