#!/usr/bin/env bash
# Set the GALAXY_TOKEN Actions secret on every ansible role repo.
#
# GitHub user accounts have no account-level secrets, so each role repo holds
# its own copy of the Galaxy API token. Regenerating the token at
# https://galaxy.ansible.com/ui/token/ revokes the old one everywhere.
#
# Usage:
#   scripts/rollout-galaxy-token.sh                    # dry run: list target repos
#   scripts/rollout-galaxy-token.sh --apply            # prompt for token, set secret
#   scripts/rollout-galaxy-token.sh --apply --rerun    # also rerun failed tag publishes
#   scripts/rollout-galaxy-token.sh --apply ansible-role-argo ansible-role-qq
#
# The token is read from $GALAXY_TOKEN, or prompted for without echo. It is
# checked against the Galaxy API before anything is written, and passed to gh
# on stdin so it never shows up in process args.
set -euo pipefail

owner="${OWNER:-bradfordwagner}"
secret_name="GALAXY_TOKEN"
# repos named ansible-role-* plus the role template
repo_pattern='^(ansible-role-.+|gh-template-ansible-role)$'

apply=false
rerun=false
repos=()
for arg in "$@"; do
  case "$arg" in
    --apply) apply=true ;;
    --rerun) rerun=true ;;
    -h|--help) sed -n '2,16p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    -*) echo "unknown flag: $arg" >&2; exit 2 ;;
    *) repos+=("$arg") ;;
  esac
done

if [ ${#repos[@]} -eq 0 ]; then
  while IFS= read -r r; do repos+=("$r"); done < <(
    gh repo list "$owner" -L 1000 --no-archived --source --json name -q '.[].name' |
      grep -E "$repo_pattern" | sort
  )
fi

if [ ${#repos[@]} -eq 0 ]; then
  echo "no repos matched $repo_pattern under $owner" >&2
  exit 1
fi

echo "target repos (${#repos[@]}):"
for r in "${repos[@]}"; do
  updated=$(gh api "repos/$owner/$r/actions/secrets/$secret_name" -q .updated_at 2>/dev/null || echo "missing")
  printf '  %-40s %s=%s\n' "$r" "$secret_name" "$updated"
done

if ! $apply; then
  echo
  echo "dry run. rerun with --apply to set $secret_name."
  exit 0
fi

token="${GALAXY_TOKEN:-}"
if [ -z "$token" ]; then
  read -rsp "galaxy api token: " token
  echo
fi
[ -n "$token" ] || { echo "empty token" >&2; exit 1; }

# same endpoint ansible-galaxy hits first: 200 valid, 401 bad token
code=$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Token $token" https://galaxy.ansible.com/api/)
if [ "$code" != "200" ]; then
  echo "galaxy rejected token (HTTP $code), nothing written" >&2
  exit 1
fi
echo "token valid"

failed=()
for r in "${repos[@]}"; do
  if printf '%s' "$token" | gh secret set "$secret_name" -R "$owner/$r" >/dev/null; then
    echo "  set   $r"
  else
    echo "  FAIL  $r"
    failed+=("$r")
    continue
  fi

  if $rerun; then
    run=$(gh run list -R "$owner/$r" -w container-tags -L 1 --json databaseId,conclusion \
      -q '.[0] | select(.conclusion == "failure") | .databaseId')
    if [ -n "$run" ]; then
      gh run rerun "$run" -R "$owner/$r" --failed && echo "        rerun container-tags run $run"
    fi
  fi
done

if [ ${#failed[@]} -gt 0 ]; then
  echo "failed: ${failed[*]}" >&2
  exit 1
fi
