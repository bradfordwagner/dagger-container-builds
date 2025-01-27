package lib

import (
	"context"
	"dagger/container-builds/internal/dagger"
	"strings"
)

type ManifestToolArgs struct {
	// base args
	Actor *dagger.Secret
	Token *dagger.Secret

	// image-tag args - ie for different os
}

func ManifestTool(
	ctx context.Context,
	container *dagger.Container,
	actor *dagger.Secret,
	token *dagger.Secret,
	image string,
	arches []string,
) (s string, err error) {
	//username, _ := actor.Plaintext(ctx)
	password, _ := token.Plaintext(ctx)
	c := container.From("mplatform/manifest-tool:alpine-v2.1.9").
		WithExec([]string{
			"manifest-tool",
			"--username", "doesntmatter",
			//"--username", fmt.Sprintf("%s", username),
			"--password", password,
			"push", "from-args",
			"--platforms", strings.Join(arches, ","),
			"--template", image + "-OS_ARCH",
			"--target", image,
		})
	return ContainerOutput(ctx, c)
}
