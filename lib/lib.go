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
