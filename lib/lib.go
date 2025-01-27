package lib

import (
	"context"
	"dagger/container-builds/internal/dagger"
	"strings"
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
		"linux/arm64": "ubuntu-24.04-arm",
		//"linux/arm64": "arm64",
	}

	// default to ubuntu-latest
	s, ok := archs[arch]
	if !ok {
		return "ubuntu-latest"
	}
	return
}

// ArchImageName returns the image name for the given image and architecture
func ArchImageName(image, arch string) (s string) {
	arch = strings.ReplaceAll(arch, "/", "_")
	return image + "-" + arch
}
