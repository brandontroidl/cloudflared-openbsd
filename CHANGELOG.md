# Changelog

All notable changes to this repository are documented here. The format is based
on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

This repository patches [`cloudflared`](https://github.com/cloudflare/cloudflared);
"upstream" below refers to the cloudflared version that `build.sh` targets by
default.

## [Unreleased]

### Added

- `contrib/cloudflared.rc`, an OpenBSD `rc.d` service script, with README
  instructions for running the binary as a managed daemon under a dedicated
  `_cloudflared` user.
- README section positioning this cross-build against the ports-tree port
  ([openbsd-port-cloudflared](https://github.com/ivoronin/openbsd-port-cloudflared)).

## [0.1.0] - 2026-06-18

### Added

- OpenBSD `SystemCollectorImpl` (`system_collector_openbsd.go`) reading memory
  (`hw.physmem`) and file-descriptor counts (`kern.maxfiles`, `kern.nfiles`)
  via `sysctl`. Adapted from the FreeBSD ports patch.
- `build.sh`, which cross-builds a static `openbsd/amd64` cloudflared, adding
  `openbsd` to the `//go:build` constraint of
  `diagnostic/network/collector_unix.go` at build time. Default upstream target:
  cloudflared 2026.6.0.
- `SECURITY.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`.

### Known limitations

- Current (used) memory is reported as `0`; OpenBSD exposes no single sysctl
  scalar for free memory, so only total memory is reported.
