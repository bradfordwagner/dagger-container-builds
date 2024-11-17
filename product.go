package main

import (
	"context"
	"dagger/container-builds/internal/dagger"
	"dagger/container-builds/lib/flavors"
	"errors"
)

// ProductJson returns the product of the config yaml in order to matrix the build
func (m *ContainerBuilds) ProductJson(
	ctx context.Context,
	src dagger.Directory,
	// +default="latest"
	version string,
) (s string, err error) {
	c, err := loadConfig(ctx, src)
	if err != nil {
		return
	}

	switch c.Flavor {
	case flavors.FlavorMirror:
		s = `{"flavor":"mirror"}`
	default:
		err = errors.New("unknown flavor")
	}

	return s, nil
}
