package mirror

import (
	"context"
	"dagger/container-builds/internal/dagger"
	"dagger/container-builds/lib"
)

// Manifest - configures manifest in registry. not meant to be run locally
func Manifest(
	ctx context.Context,
	container *dagger.Container,
	configString string,
	version string,
	actor *dagger.Secret,
	token *dagger.Secret,
) (o string, err error) {
	c, err := loadConfig(configString)
	if err != nil {
		return
	}

	// run manifest tool for each build
	for _, b := range c.Builds {
		image := imageTag(c, b, version)
		o, err = lib.ManifestTool(ctx, container, actor, token, image, b.Architectures)
		if err != nil {
			return
		}
	}

	return
}
