package mirror

import (
	"dagger/container-builds/lib"
	"encoding/json"
	"fmt"
)

type ProductFormat struct {
	Index        int    `json:"index"`
	Repo         string `json:"repo"`
	Tag          string `json:"tag"`
	Architecture string `json:"arch"`
	Runner       string `json:"runner"`
	TargetImage  string `json:"target_image"` // without architecture suffix
}

func Product(
	configString string,
	version string,
) (products []ProductFormat, err error) {
	c, err := loadConfig(configString)
	if err != nil {
		return
	}

	// create a list of products
	var i int
	for _, b := range c.Builds {
		for _, a := range b.Architectures {
			runner := lib.ArchToRunner(a)
			products = append(products, ProductFormat{
				Architecture: a,
				Index:        i,
				Repo:         b.Repo,
				Runner:       runner,
				Tag:          b.Tag,
				TargetImage:  imageTag(c, b, version),
			})
			i++
		}
	}

	return
}

func imageTag(c Config, b Build, version string) string {
	repo := b.Repo
	if b.RepoOverride != "" {
		repo = b.RepoOverride
	}
	return fmt.Sprintf("%s:%s-%s_%s", c.TargetRepo, version, repo, b.Tag)
}

// ProductJson returns the cartesian product of all builds as a json string, used for github actions matrix
func ProductJson(
	configString string,
	version string,
) (o string, err error) {
	products, err := Product(configString, version)
	bytes, err := json.Marshal(products)
	if err != nil {
		return
	}
	o = string(bytes)
	return
}
