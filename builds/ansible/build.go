package ansible

import (
	"context"
	"dagger/container-builds/internal/dagger"
	"dagger/container-builds/lib"
	"dagger/container-builds/lib/flavors"
	"encoding/json"
	"strings"
)

// Build - builds the container image
func BuildContainer(
	ctx context.Context,
	client *dagger.Client,
	src dagger.Directory,
	index int,
	version string,
	configString string,
) (o string, err error) {
	products, err := Product(configString, version)
	if err != nil {
		return
	}

	config, err := loadConfig(configString)
	if err != nil {
		return "", err
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

	// build the container
	container := client.Container(dagger.ContainerOpts{
		Platform: dagger.Platform(product.Architecture),
	}).From(product.UpstreamImage).WithDirectory("/src", &src)
	//container = dag.Lib().InvalidateCache(invalidateCache, container)

	// find requirements
	requirements := []string{"requirements.yml", "meta/requirements.yml"}
	for _, requirement := range requirements {
		var exists bool
		if exists, err = fileExists(ctx, src, requirement); err != nil {
			return
		}
		if exists {
			container, err = container.WithExec([]string{"ansible-galaxy", "install", "-r", requirement}).Sync(ctx)
			if err != nil {
				return
			}
		}
	}

	// run playbook
	playbooks := []string{"test.yml", "playbook.yml"}
	for _, playbook := range playbooks {
		var exists bool
		if exists, err = fileExists(ctx, src, playbook); err != nil {
			return
		}
		if exists {
			container, err = container.WithExec([]string{"ansible-playbook", playbook}).Sync(ctx)
			if err != nil {
				return
			}
		}
	}
	o = strings.Join([]string{target, productJson}, "\n")

	if config.Flavor == flavors.FlavorAnsiblePlaybook {
		_, err = container.Publish(ctx, target)
	}

	return
}

// fileExists reports whether path is a regular file in src
func fileExists(ctx context.Context, src dagger.Directory, path string) (bool, error) {
	return src.Exists(ctx, path, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeRegularType})
}
