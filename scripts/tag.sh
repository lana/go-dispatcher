#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Run unit tests
set -eufo pipefail
export SHELLOPTS
IFS=$'\t\n'

command -v git >/dev/null 2>&1 || { echo 'please install git or use image that has it'; exit 1; }

git fetch --tags

# TODO: Consider something more robust, like piping all tags through `grep`, so we can use regular expression matches, rather than fnmatch.
# get latest tag that looks like a semver (with or without v)
tag=$(git for-each-ref --sort=-v:refname --count=1 --format '%(refname)' refs/tags/[0-9]*.[0-9]*.[0-9]* refs/tags/v[0-9]*.[0-9]*.[0-9]* | cut -d / -f 3- | sed 's/v//g')

if [ -z "$tag" ]; then
    new="0.0.1"
else
    ancestor=$(git merge-base HEAD^1 refs/remotes/origin/master)
    log=$(git log "$ancestor"..HEAD | grep '^\s')
    case "$log" in
        *#major*) new=$(semver-cli inc major "${tag}");;
        *#minor*) new=$(semver-cli inc minor "${tag}");;
        *#patch*) new=$(semver-cli inc patch "${tag}");;
        *)
            echo "At least one commit description must contain: #major, #minor, or #patch"
            exit 1
        ;;
    esac
fi

if [ -z "$new" ]; then
    echo "Failed to generate new tag. This is a bug."
    exit 1
fi

# TODO: Switch this to GitLab
#git remote add org "https://oauth2:${GITLAB_ACCESS_TOKEN}@gitlab.lana.tech/${CI_PROJECT_DIR}"
git remote add org "https://lana-dev:${GITHUB_ACCESS_TOKEN}@github.com/lana/go-dispatcher"

git tag v"$new"
git push org v"$new"
