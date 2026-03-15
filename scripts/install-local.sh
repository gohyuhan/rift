#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="rift"

echo "Building rift..."
go build -o "$BINARY_NAME" .

echo "Installing to $INSTALL_DIR/$BINARY_NAME (may require sudo)..."
if [ -w "$INSTALL_DIR" ]; then
    mv -f "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mv -f "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

echo "Done: $(which rift)"
rift --version 2>/dev/null || true
