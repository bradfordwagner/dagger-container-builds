package mirror

import (
	"context"
	"dagger/container-builds/internal/dagger"
	"dagger/container-builds/lib"
	"encoding/json"
	"fmt"
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

	//dockerfile setup
	var runSteps string
	if len(product.Commands) > 0 {
		scriptEntries, _ := src.Directory("scripts").Entries(ctx)
		if len(scriptEntries) > 0 {
			runSteps += "COPY scripts /tmp/build/scripts\nRUN chmod +x /tmp/build/scripts/*.sh\n"
		}
		for _, cmd := range product.Commands {
			runSteps += fmt.Sprintf("RUN %s\n", cmd)
		}
		if len(scriptEntries) > 0 {
			runSteps += "RUN rm -rf /tmp/build/scripts\n"
		}
	}
	dockerfile := fmt.Sprintf("FROM %s:%s\n%s", product.Repo, product.Tag, runSteps)
	d := src.WithNewFile("Dockerfile", dockerfile)
	container := d.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Platform: dagger.Platform(product.Architecture),
	})
	o = strings.Join([]string{target, productJson, dockerfile}, "\n")

	// publish only through pipeline
	_, err = container.Publish(ctx, target)

	return
}
