#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(pwd)"
CLI_DIR="$ROOT_DIR/cli"

echo "Running TrustSSH CLI unit tests"

cd "$CLI_DIR"
go test -v ./...

echo
echo "All Go tests passed"
