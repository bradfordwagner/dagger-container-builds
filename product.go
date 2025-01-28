package main

import (
	"context"
	"dagger/container-builds/builds/ansible"
	"dagger/container-builds/builds/custom"
	"dagger/container-builds/builds/mirror"
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
	c, yml, err := loadConfig(ctx, src)
	if err != nil {
		return
	}

	switch c.Flavor {
	case flavors.FlavorAnsibleRole:
		s, err = ansible.ProductJson(yml, version)
	case flavors.FlavorAnsiblePlaybook:
		s, err = ansible.ProductJson(yml, version)
	case flavors.FlavorCustom:
		s, err = custom.ProductJson(yml, version)
	case flavors.FlavorMirror:
		s, err = mirror.ProductJson(yml, version)
	default:
		err = errors.New("unknown flavor")
	}

	return
}
