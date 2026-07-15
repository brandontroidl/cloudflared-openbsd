# cloudflared on OpenBSD

A patch and build script that make Cloudflare's [`cloudflared`](https://github.com/cloudflare/cloudflared)
compile and run on OpenBSD.

Upstream `cloudflared` ships a diagnostic system collector for Linux, macOS, and
FreeBSD, but not OpenBSD, so recent releases fail to build for `openbsd/amd64`.
This repository supplies the missing collector and a script that cross-builds a
static binary from an unmodified upstream tree.

The tunnel data path is upstream code, untouched. The added collector runs only
under `cloudflared tunnel diag`.

## Requirements

- Go (a recent toolchain; the build uses `GOTOOLCHAIN=auto`)
- `git`
- Any OS that can cross-compile for OpenBSD (the build host does not need to be
  OpenBSD)

## Build

```sh
sh build.sh            # builds against cloudflared 2026.6.0 (default)
sh build.sh 2026.5.0   # or pin a specific upstream version
```

`build.sh` clones `cloudflared` at the requested tag, applies the two changes
described below, and builds with `CGO_ENABLED=0`. The result is a static binary
named `cloudflared-openbsd-amd64`.

## Install

On the OpenBSD host:

```sh
doas install -o root -g bin -m 755 cloudflared-openbsd-amd64 /usr/local/sbin/cloudflared
cloudflared --version
```

## What the patch changes

Two files, both confined to the diagnostic collector:

- **`system_collector_openbsd.go`** (added) implements `SystemCollectorImpl` for
  OpenBSD. It reads memory and file-descriptor counts through `sysctl`:
  `hw.physmem` for total memory, `kern.maxfiles` and `kern.nfiles` for
  descriptors. Adapted from the FreeBSD ports patch.
- **`diagnostic/network/collector_unix.go`** (modified at build time) has its
  `//go:build` constraint rewritten to include `openbsd`, so the shared Unix
  network collector compiles for the platform. `build.sh` performs this rewrite;
  the upstream file is not edited in this repository.

## Limitations

- Current (used) memory is reported as `0`. OpenBSD exposes no single sysctl
  scalar for free memory, so the collector reports total memory only. This
  affects the `tunnel diag` output, nothing else.
- Only `openbsd/amd64` is built. Other architectures are untested.

## Relationship to upstream

This is a bridge until OpenBSD support lands in `cloudflared` itself. If you run
OpenBSD, consider carrying these changes upstream so the port is maintained
there rather than here.

## Contributing, security, conduct

- Reporting a vulnerability: [`SECURITY.md`](SECURITY.md)
- Submitting changes: [`CONTRIBUTING.md`](CONTRIBUTING.md)
- Community expectations: [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md)
- Version history: [`CHANGELOG.md`](CHANGELOG.md)

## License

Apache-2.0, matching `cloudflared`. See [`LICENSE`](LICENSE).

## Trademark and affiliation

This is not an official Cloudflare project and is not affiliated with or endorsed
by Cloudflare, Inc. "Cloudflare" and "cloudflared" are trademarks of Cloudflare,
Inc., used here only to identify the upstream software.
