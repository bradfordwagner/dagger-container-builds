package main

import (
	"context"
	"dagger/container-builds/builds/mirror"
	"dagger/container-builds/internal/dagger"
	"dagger/container-builds/lib/flavors"
	"errors"
)

func (m *ContainerBuilds) Manifest(
	ctx context.Context,
	src dagger.Directory,
	// +default="latest"
	version string,
	// GitHub actor, --token=env:github_actor,--token=cmd:"gh auth token"
	actor *dagger.Secret,
	// GitHub API token, --token=env:github_token,--token=cmd:"gh auth token"
	token *dagger.Secret,
) (s string, err error) {
	c, yml, err := loadConfig(ctx, src)
	if err != nil {
		return
	}

	container := dag.Container()
	switch c.Flavor {
	case flavors.FlavorMirror:
		s, err = mirror.Manifest(ctx, container, yml, version, actor, token)
	default:
		err = errors.New("unknown flavor")
	}

	return
}
