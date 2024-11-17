package main

import (
	"context"
	"dagger/container-builds/internal/dagger"
	"dagger/container-builds/lib"
	"dagger/container-builds/lib/flavors"

	"gopkg.in/yaml.v2"
)

type ContainerBuilds struct{}

type BaseConfig struct {
	Flavor flavors.Flavor `yaml:"flavor"`
}

// Returns a container that echoes whatever string argument is provided
func (m *ContainerBuilds) Echo() string {
	return "hello world"
}

func (m *ContainerBuilds) Config(
	ctx context.Context,
	src dagger.Directory,
) (s string, err error) {
	return lib.OpenConfigYaml(ctx, src)
}

// loadConfig loads the config.yaml from the source directory
func loadConfig(ctx context.Context, src dagger.Directory) (c BaseConfig, err error) {
	yml, _ := lib.OpenConfigYaml(ctx, src)
	err = yaml.Unmarshal([]byte(yml), &c)
	return
}
