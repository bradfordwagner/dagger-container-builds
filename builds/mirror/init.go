package mirror

import (
	"context"

	"gopkg.in/yaml.v2"
)

// Init creates an example yaml config for cicd to use
func Init(
	ctx context.Context,
) (s string, err error) {
	// default config
	c := Config{
		Flavor:     "mirror",
		TargetRepo: "ghcr.io/bradfordwagner/template-mirror",
		Builds: []Build{
			{Repo: "alpine", Tag: "3.18", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "alpine", Tag: "3.19", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "alpine", Tag: "3.20", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "archlinux", Tag: "latest", Architectures: []string{"linux/amd64"}},
			{Repo: "debian", Tag: "bookworm", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "debian", Tag: "bookworm-slim", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "rockylinux", Tag: "8", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "rockylinux", Tag: "9", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "ubuntu", Tag: "jammy", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "ubuntu", Tag: "mantic", Architectures: []string{"linux/amd64", "linux/arm64"}},
			{Repo: "ubuntu", Tag: "noble", Architectures: []string{"linux/amd64", "linux/arm64"}},
		},
	}

	// convert to yaml
	b, err := yaml.Marshal(c)
	if err != nil {
		return
	}

	// return yaml
	return string(b), nil
}
