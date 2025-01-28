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
		"linux/amd64": "ubuntu-24.04",
	}

	// default to ubuntu-latest
	s, ok := archs[arch]
	if !ok {
		return "ubuntu-24.04"
	}
	return
}

// ArchImageName returns the image name for the given image and architecture
func ArchImageName(image, arch string) (s string) {
	arch = strings.ReplaceAll(arch, "/", "_")
	return image + "-" + arch
}

func FileContents(ctx context.Context, dir *dagger.Directory, path string) (contents string, err error) {
	if file := dir.File(path); file != nil {
		contents, err = file.Contents(ctx)
	}
	return
}
