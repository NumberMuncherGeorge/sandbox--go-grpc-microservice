#!/usr/bin/env bash
set -euo pipefail

# Download googleapis for HTTP annotations if not already present
PROTO_DEPS_DIR="third_party/googleapis"

if [ ! -d "$PROTO_DEPS_DIR" ]; then
    echo "Downloading googleapis for HTTP annotations..."
    mkdir -p third_party
    git clone --depth 1 --branch master https://github.com/googleapis/googleapis.git "$PROTO_DEPS_DIR"
    echo "✓ googleapis downloaded to $PROTO_DEPS_DIR"
else
    echo "✓ googleapis already exists at $PROTO_DEPS_DIR"
fi
