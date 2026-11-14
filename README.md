# dagger-container-builds

> **Deprecated.** Superseded by [bradfordwagner/gh-template-lib](https://github.com/bradfordwagner/gh-template-lib),
> reusable GitHub Actions workflows that build with plain `docker buildx` instead of Dagger. No new releases;
> `0.4.1` is the last tag and stays in place for anything still pinned to it.

## Migrating

Replace the inline Dagger jobs in `container_branches.yml` and `container_tags.yml` with the single `container`
job from the matching template's latest tag. `config.yaml` is unchanged: gh-template-lib reads the same schema and
emits the same product matrix and image tags.

| `flavor` | Template | Reusable workflow |
|---|---|---|
| `mirror` | [gh-template-mirror-image](https://github.com/bradfordwagner/gh-template-mirror-image) | `gh-template-lib/.github/workflows/mirror.yml` |
| `custom` | [gh-template-custom-image](https://github.com/bradfordwagner/gh-template-custom-image) | `gh-template-lib/.github/workflows/custom.yml` |
| `ansible_playbook` | [gh-template-ansible-playbook](https://github.com/bradfordwagner/gh-template-ansible-playbook) | `gh-template-lib/.github/workflows/ansible_playbook.yml` |
| `ansible_role` | [gh-template-ansible-role](https://github.com/bradfordwagner/gh-template-ansible-role) | `gh-template-lib/.github/workflows/ansible_role.yml` |

See gh-template-lib's [migration notes](https://github.com/bradfordwagner/gh-template-lib/blob/main/docs/migration.md)
for the behaviour differences.
