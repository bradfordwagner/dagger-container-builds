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

// loadConfig loads the config.yaml from the source directory
func loadConfig(ctx context.Context, src dagger.Directory) (c BaseConfig, yml string, err error) {
	yml, _ = lib.OpenConfigYaml(ctx, src)
	err = yaml.Unmarshal([]byte(yml), &c)
	return
}
