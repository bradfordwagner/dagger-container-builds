package ansible

import (
	"dagger/container-builds/lib/flavors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Flavor     flavors.Flavor `yaml:"flavor"`
	TargetRepo string         `yaml:"target_repo"`
	Builds     []Build        `yaml:"builds"`
	Upstream   Upstream       `yaml:"upstream"`
}

type Upstream struct {
	Repo string `yaml:"repo"`
	Tag  string `yaml:"tag"`
}

type Build struct {
	OS            string   `yaml:"os"`
	Architectures []string `yaml:"archs"`
}

func loadConfig(configString string) (c Config, err error) {
	err = yaml.Unmarshal([]byte(configString), &c)
	return
}
