#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Run unit tests
set -eufo pipefail
export SHELLOPTS
IFS=$'\t\n'

command -v git >/dev/null 2>&1 || { echo 'please install git or use image that has it'; exit 1; }

git fetch --tags

ancestor=$(git merge-base HEAD refs/remotes/origin/master)
log=$(git log "$ancestor"..HEAD | grep '^\s')
case "$log" in
    *#major*|*#minor*|*#patch*) # All Good
    ;;
    *)
        echo "At least one commit description must contain: #major, #minor, or #patch"
        echo "See https://gitlab.lana.tech/eng/help/-/wikis/semver.md for details"
        exit 1
    ;;
esac
