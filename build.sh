#!/bin/bash

export WORKINGDIR=$(pwd)

target=${1:-build}
arch=$(go env GOHOSTARCH)
platform=$(go env GOHOSTOS)
binPath="$(pwd)/dist/prm_${platform}_${arch}"
binPath2="$(pwd)/dist/notel_prm_${platform}_${arch}"

if [ "$target" == "build" ]; then
  # Set goreleaser to build for current platform only
  if [ -z "${HONEYCOMB_API_KEY}" ]; then
    export HONEYCOMB_API_KEY="not_set"
  fi
  if [ -z "${HONEYCOMB_DATASET}" ]; then
    export HONEYCOMB_DATASET="not_set"
  fi
  goreleaser build --snapshot --rm-dist --single-target
  git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/fantastic-adventure "$binPath/prm-tools"
  cp -r "$binPath/prm-tools" "$binPath2/prm-tools"
elif [ "$target" == "quick" ]; then
  go build -o ${binPath}/prm -tags telemetry
  go build -o ${binPath2}/prm
elif [ "$target" == "package" ]; then
  git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/fantastic-adventure "prm-tools"
  goreleaser --skip-publish --snapshot --rm-dist
fi
