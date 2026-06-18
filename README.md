# cloudflared on OpenBSD

cloudflared's diagnostic collector has no OpenBSD implementation, so current
releases don't build on OpenBSD. This adds the missing collector and a script to
cross-build a static `openbsd/amd64` binary.

Only `cloudflared tunnel diag` touches this code. The tunnel itself is unchanged.

Not an official Cloudflare project and not affiliated with Cloudflare.
"Cloudflare" and "cloudflared" are trademarks of Cloudflare, Inc.

## Build

Needs Go. `build.sh` clones cloudflared at the given version, applies the two
changes below, and builds with `CGO_ENABLED=0`.

```sh
sh build.sh            # cloudflared 2026.6.0
sh build.sh 2026.5.0   # or pick a version
```

Output is `cloudflared-openbsd-amd64`.

## Install

On the OpenBSD host:

```sh
doas install -o root -g bin -m 755 cloudflared-openbsd-amd64 /usr/local/sbin/cloudflared
cloudflared --version
```

## The changes

`system_collector_openbsd.go` is a new `SystemCollectorImpl` for OpenBSD. Memory
and file-descriptor numbers come from sysctl: `hw.physmem` for memory,
`kern.maxfiles` and `kern.nfiles` for descriptors. Based on the FreeBSD ports
patch.

`diagnostic/network/collector_unix.go` excludes OpenBSD in its `//go:build`
constraint. `build.sh` rewrites that line to add `openbsd` so the unix collector
compiles.

## License

Apache-2.0, same as cloudflared (see `LICENSE`). Files changed from upstream:
`system_collector_openbsd.go` is added, and the build constraint in
`diagnostic/network/collector_unix.go` is modified at build time.
