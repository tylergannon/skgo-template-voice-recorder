#!/bin/sh
# Build for the exact invoking SKGo host using its supported shared helper.
set -eu
repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_root"
go run ./template/stage
exec go run github.com/tylergannon/skgo/cmd/skgo-plugin-build \
  --source "$repo_root" --package ./template "$@"
