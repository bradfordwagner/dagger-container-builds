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

// Build kicks of the container build
// --src=. \\
// --index={{ `\${{ matrix.product.index }}` }} \\
// --version=\${version} \\
// --isDev=false

func (m *ContainerBuilds) Build(
	ctx context.Context,
	src dagger.Directory,
	// +default=0
	index int,
	// +default="latest"
	version string,
) (s string, err error) {
	c, yml, err := loadConfig(ctx, src)
	if err != nil {
		return
	}

	client := dagger.Connect()
	switch c.Flavor {
	case flavors.FlavorAnsibleRole:
		s, err = ansible.BuildContainer(ctx, client, src, index, version, yml)
	case flavors.FlavorAnsiblePlaybook:
		s, err = ansible.BuildContainer(ctx, client, src, index, version, yml)
	case flavors.FlavorCustom:
		s, err = custom.BuildContainer(ctx, src, index, version, yml)
	case flavors.FlavorMirror:
		s, err = mirror.BuildContainer(ctx, src, index, version, yml)
	default:
		err = errors.New("unknown flavor")
	}

	return
}
