#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
socdir="$workspace/src/github.com/allsportschain"
if [ ! -L "$socdir/go-allsportschain" ]; then
    mkdir -p "$socdir"
    cd "$socdir"
    ln -s ../../../../../. go-allsportschain
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$socdir/go-allsportschain"
PWD="$socdir/go-allsportschain"

# Launch the arguments with the configured environment.
exec "$@"
