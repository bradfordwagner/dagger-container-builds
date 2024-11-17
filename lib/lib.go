package lib

import (
	"context"
	"dagger/container-builds/internal/dagger"
)

// Returns lines that match a pattern in the files of the provided Directory
func OpenConfigYaml(
	ctx context.Context,
	src dagger.Directory,
) (s string, err error) {
	configFile := src.File("config.yaml")
	return configFile.Contents(ctx)
}
