package lib

import (
	"context"
	"dagger/container-builds/internal/dagger"
)

// ContainerOutput returns the output of a container as a string if stderr exists return that, else stdout
func ContainerOutput(ctx context.Context, c *dagger.Container) (s string, err error) {
	s, err = c.Stderr(ctx)
	if err != nil || s != "" {
		return
	}
	return c.Stdout(ctx)
}
