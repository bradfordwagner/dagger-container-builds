# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Dagger module (Go SDK, engine `v0.15.2`, module name `container_builds`) that builds and publishes multi-arch container images from a `config.yaml` in a consumer repo. It is driven from GitHub Actions: one call produces a build matrix, one call per matrix entry builds a single arch image, and a final call stitches the per-arch images into a multi-arch manifest.

## Commands

Generated code (`dagger.gen.go`, `internal/dagger`, `internal/querybuilder`, `internal/telemetry`) is gitignored. Regenerate it before `go build`/`go vet` or IDE tooling will resolve imports:

```sh
dagger develop          # regenerate SDK bindings
go build ./...
go vet ./...
```

There are no tests. Exercise functions through the Dagger CLI against a directory holding a `config.yaml`:

```sh
dagger functions
dagger call init --flavor=mirror
dagger call product-json --src=<dir> --version=1.0.0
dagger call build --src=<dir> --index=0 --version=1.0.0      # publishes to the registry
dagger call manifest --src=<dir> --version=1.0.0 --actor=env:GITHUB_ACTOR --token=cmd:"gh auth token"
```

`build` and `manifest` push to real registries; they are intended to run in CI, not locally.

## Architecture

Top-level files (`init.go`, `product.go`, `build.go`, `manifest.go`) each define one exported method on `ContainerBuilds`. Each reads `config.yaml` from `src` via `loadConfig` (`main.go`), switches on the `flavor` field, and delegates to a package under `builds/<flavor>/`. Each flavor package mirrors the same file layout: `config.go` (YAML schema + `loadConfig`), `init.go` (example config), `product.go`, `build.go`, `manifest.go`.

Flavors (`lib/flavors`): `mirror`, `custom`, `ansible_role`, `ansible_playbook`. The two ansible flavors share `builds/ansible`.

Pipeline contract:
- **Product** expands `builds × archs` into a flat list of `ProductFormat`, each with a sequential `index`, a GitHub runner (`lib.ArchToRunner`: arm64 maps to `ubuntu-24.04-arm`, everything else to `ubuntu-24.04`), and a `target_image` without the arch suffix. `ProductJson` emits this as the Actions matrix.
- **Build** recomputes the product list and picks entry `index`. It pushes to `target_image` plus an arch suffix (`lib.ArchImageName`, e.g. `...-linux_amd64`).
- **Manifest** runs `mplatform/manifest-tool` (`lib.ManifestTool`) per build, combining the `-OS_ARCH` images into `target_image`. Only the token is used as the password; the username is a placeholder.

The product ordering must be deterministic, because `index` from the matrix job is resolved again inside `Build`.

Per-flavor build behavior:
- `mirror`: generates a Dockerfile `FROM repo:tag`, adding `RUN` lines from `commands`. If `src/scripts` exists it is copied to `/src/scripts`, made executable, and removed afterwards. `repo_override` renames the repo in the target tag (`<target_repo>:<version>-<repo>_<tag>`).
- `custom`: runs `DockerBuild` on the consumer's own Dockerfile, passing `build.args` plus `OS=<upstream.repo>:<upstream.tag>-<os>` as build args. Target tag is `<target_repo>:<version>-<os>`.
- `ansible_*`: starts from `<upstream.repo>:<upstream.tag>-<os>` and mounts `src` at `/src`. It runs `ansible-galaxy install` for `requirements.yml`/`meta/requirements.yml` if present, then `ansible-playbook` for `test.yml`/`playbook.yml` if present. Only `ansible_playbook` publishes; `ansible_role` is a test-only build. This path uses a `dagger.Connect()` client rather than `dag`.

## Flavor template repos

Each flavor has a template repo, cloned under `/Users/bwagner/workspace/github/bradfordwagner/gh_templates/`. To iterate on a flavor, change the matching template first: its `config.yaml` and `.github/workflows/container_{branches,tags}.yml` are the real consumers of this module. Keep this table current when adding a flavor or template.

| Flavor | Local dir (under `gh_templates/`) | GitHub repo | Workflows call |
|---|---|---|---|
| `mirror` | `bradfordwagner.template.container.mirror.image` | `bradfordwagner/gh-template-mirror-image` | product-json, build, manifest |
| `ansible_playbook` | `bradfordwagner.template.ansible.playbook` | `bradfordwagner/gh-template-ansible-playbook` | product-json, build, manifest |
| `ansible_role` | `bradfordwagner.template.ansible.role` | `bradfordwagner/gh-template-ansible-role` | product-json, build (no manifest; publishes to Galaxy on tags) |
| `custom` | `bradfordwagner.template.container.custom.image` | `bradfordwagner/gh-template-custom-image` | product-json, build, manifest |

Real `custom` consumer: `containers/custom/bradfordwagner.container.custom.ansible` (`bradfordwagner/container-ansible`, builds `ghcr.io/bradfordwagner/ansible`). The template was derived from it.

Template clones are managed by `~/dotfiles/tasks/git-modules-bradfordwagner.yml` (the `templates_dir` loop). Add new templates there too.

Templates pin the module by tag (`dagger -m https://github.com/bradfordwagner/dagger-container-builds@<tag>`) with Dagger CLI `0.15.2`. Pins are currently `@0.2.0` in mirror and ansible, and `@0.3.0` in custom. A change here reaches templates only after you cut a new tag and bump that pin. To test an untagged change, point the template's `-m` at a branch or commit.

Other repos in `gh_templates/` (go.cli, go.lib, helm.chart, tf.module, tf.workspaces) do not use this module.

## Known gaps

Known gaps when adding to a flavor: `Init` only handles `mirror`, and `Manifest` has no `ansible_role` case. Adding a flavor means updating all four top-level switches plus `lib/flavors`.
