#!/usr/bin/env bash
set -e

# Check protoc is installed
if ! command -v protoc &>/dev/null; then
  echo "protoc not found. Install with: brew install protobuf"
  exit 1
fi

# Check protoc-gen-go is installed
if ! command -v protoc-gen-go &>/dev/null; then
  echo "protoc-gen-go not found. Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
  exit 1
fi

# Compile all .proto files
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  proto/*.proto

# Sync dependencies
go mod tidy

echo "Proto generation done."
