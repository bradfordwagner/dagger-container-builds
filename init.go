package main

import (
	"context"
	"dagger/container-builds/builds/mirror"
	"dagger/container-builds/lib/flavors"
	"errors"
)

func (m *ContainerBuilds) Init(
	ctx context.Context,
	flavor string,
) (s string, err error) {
	switch flavors.FromString(flavor) {
	case flavors.FlavorMirror:
		s, err = mirror.Init(ctx)
	default:
		err = errors.New("unknown flavor")
	}

	return
}
