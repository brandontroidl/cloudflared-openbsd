#!/bin/sh
# Cross-builds a static openbsd/amd64 cloudflared with the OpenBSD diagnostic
# collector added. Clones cloudflared at the given version (default 2026.6.0),
# drops in system_collector_openbsd.go, adds openbsd to the //go:build line of
# diagnostic/network/collector_unix.go, and builds with CGO_ENABLED=0.
# Output: cloudflared-openbsd-amd64
#
# Usage: sh build.sh [version]
set -e
VER="${1:-2026.6.0}"
HERE=$(cd "$(dirname "$0")" && pwd)

rm -rf cloudflared
git clone --depth 1 --branch "$VER" https://github.com/cloudflare/cloudflared.git
cd cloudflared

# patch 1 — add openbsd to the diagnostic network collector build tag
f=diagnostic/network/collector_unix.go
sed '1s#.*#//go:build darwin || linux || openbsd#' "$f" > "$f.tmp" && mv "$f.tmp" "$f"

# patch 2 — OpenBSD system collector
cp "$HERE/system_collector_openbsd.go" diagnostic/system_collector_openbsd.go

GOTOOLCHAIN=auto GOOS=openbsd GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-mod=vendor \
  go build -ldflags "-s -w -X main.Version=$VER" \
  -o "$HERE/cloudflared-openbsd-amd64" ./cmd/cloudflared

echo "built $HERE/cloudflared-openbsd-amd64  ($VER)"
