#!/bin/sh
if [ -n "$BASH_VERSION" ]; then
    export GOROOT=$(go env GOROOT)
elif [ -n "$FISH_VERSION" ]; then
    set -x GOROOT (go env GOROOT)
fi