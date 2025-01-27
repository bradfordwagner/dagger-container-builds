package custom

import (
	"context"
	"dagger/container-builds/internal/dagger"
	"dagger/container-builds/lib"
	"encoding/json"
	"strings"
)

// Build - builds the container image
func BuildContainer(
	ctx context.Context,
	src dagger.Directory,
	// +default=0
	index int,
	// +default="latest"
	version string,
	configString string,
) (o string, err error) {
	products, err := Product(configString, version)
	if err != nil {
		return
	}

	// load product config
	product := products[index]
	b, err := json.Marshal(product)
	if err != nil {
		return
	}
	productJson := string(b)

	// set target image
	target := lib.ArchImageName(product.TargetImage, product.Architecture)

	// find build to pull args
	c, err := loadConfig(configString)
	if err != nil {
		return
	}
	var build Build
	for _, b := range c.Builds {
		if b.OS == product.OS {
			build = b
			break
		}
	}

	//dockerfile setup
	var buildArgs []dagger.BuildArg
	for k, v := range build.Args {
		buildArgs = append(buildArgs, dagger.BuildArg{Name: k, Value: v})
	}
	buildArgs = append(buildArgs, dagger.BuildArg{Name: "OS", Value: product.UpstreamImage})
	container := src.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Platform:  dagger.Platform(product.Architecture),
		BuildArgs: buildArgs,
	})

	o = strings.Join([]string{target, productJson}, "\n")

	// publish only through pipeline
	_, err = container.Publish(ctx, target)

	return
}
