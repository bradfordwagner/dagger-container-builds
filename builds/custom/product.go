package custom

import (
	"dagger/container-builds/lib"
	"encoding/json"
	"fmt"
)

type ProductFormat struct {
	Architecture  string `json:"arch"`
	OS            string `json:"os"`
	Runner        string `json:"runner"`
	TargetImage   string `json:"target_image"` // without architecture suffix
	UpstreamImage string `json:"upstream_image"`

	// required for pipeline info
	Index   int    `json:"index"`
	Display string `json:"display"`
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
				Architecture:  a,
				Index:         i,
				Runner:        runner,
				TargetImage:   imageTag(c, b, version),
				OS:            b.OS,
				UpstreamImage: fmt.Sprintf("%s:%s-%s", c.Upstream.Repo, c.Upstream.Tag, b.OS),
				Display:       fmt.Sprintf("%s-%s(%s)", b.OS, version, a),
			})
			i++
		}
	}

	return
}

func imageTag(c Config, b Build, version string) string {
	return fmt.Sprintf("%s:%s-%s", c.TargetRepo, version, b.OS)
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
