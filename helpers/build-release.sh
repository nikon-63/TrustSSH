#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CLI_DIR="$ROOT_DIR/cli"
VERSION_FILE="$ROOT_DIR/VERSION"
DIST_DIR="$ROOT_DIR/dist/release"

if [[ ! -f "$VERSION_FILE" ]]; then
    echo "ERROR: VERSION file not found at $VERSION_FILE"
    exit 1
fi

if [[ ! -d "$CLI_DIR" ]]; then
    echo "ERROR: cli directory not found at $CLI_DIR"
    exit 1
fi

VERSION="$(tr -d '\n' < "$VERSION_FILE")"

if [[ -z "$VERSION" ]]; then
    echo "ERROR: VERSION file is empty"
    exit 1
fi

LDFLAGS="-s -w -X github.com/nikon-63/TrustSSH/cli/cmd.Version=${VERSION}"

echo "Building TrustSSH release assets"
echo "Version: $VERSION"
echo

rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

build_asset() {
    local goos="$1"
    local goarch="$2"
    local os_label="$3"
    local arch_label="$4"

    local build_dir="$TMP_DIR/${goos}_${goarch}"
    local asset_name="trustssh_${os_label}_${arch_label}.tar.gz"

    mkdir -p "$build_dir"

    echo "Building $goos/$goarch..."

    (
        cd "$CLI_DIR"

        CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
        go build \
            -trimpath \
            -ldflags "$LDFLAGS" \
            -o "$build_dir/trustssh" \
            .
    )

    chmod +x "$build_dir/trustssh"

    tar -czf "$DIST_DIR/$asset_name" -C "$build_dir" trustssh

    echo "Created dist/release/$asset_name"
    echo
}

build_asset "darwin" "arm64" "Darwin" "arm64"
build_asset "darwin" "amd64" "Darwin" "x86_64"
build_asset "linux"  "arm64" "Linux"  "arm64"
build_asset "linux"  "amd64" "Linux"  "x86_64"

echo "Generating checksums..."

(
  cd "$DIST_DIR"
  shasum -a 256 *.tar.gz > trustssh_checksums.txt
)

echo
echo "Release assets created in:"
echo "  $DIST_DIR"
echo
ls -lh "$DIST_DIR"