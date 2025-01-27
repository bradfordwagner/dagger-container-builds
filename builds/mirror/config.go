package mirror

import (
	"dagger/container-builds/lib/flavors"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Flavor     flavors.Flavor `yaml:"flavor"`
	TargetRepo string         `yaml:"target_repo"`
	Builds     []Build        `yaml:"builds"`
}

type Build struct {
	Repo          string   `yaml:"repo"`
	RepoOverride  string   `yaml:"repo_override"` // renames repo to override in the target image
	Tag           string   `yaml:"tag"`
	Architectures []string `yaml:"archs"`
}

func loadConfig(configString string) (c Config, err error) {
	err = yaml.Unmarshal([]byte(configString), &c)
	return
}
