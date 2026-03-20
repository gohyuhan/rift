#!/usr/bin/env bash
set -e

# Check protoc is installed
if ! command -v protoc &>/dev/null; then
  echo "protoc not found. Install with: brew install protobuf"
  exit 1
fi

# Check protoc-gen-go is installed (also check $GOPATH/bin and $HOME/go/bin)
if ! command -v protoc-gen-go &>/dev/null; then
  GOBIN="$(go env GOPATH)/bin"
  if [ -x "$GOBIN/protoc-gen-go" ]; then
    export PATH="$PATH:$GOBIN"
  else
    echo "protoc-gen-go not found. Install with:"
    echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    echo "Then add Go's bin directory to your PATH:"
    echo "  export PATH=\"\$PATH:\$(go env GOPATH)/bin\""
    exit 1
  fi
fi

# Compile all .proto files
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  proto/*.proto

# Sync dependencies
go mod tidy

echo "Proto generation done."
