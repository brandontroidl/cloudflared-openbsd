# cloudflared on OpenBSD

A patch and build script that make Cloudflare's [`cloudflared`](https://github.com/cloudflare/cloudflared)
compile and run on OpenBSD.

Upstream `cloudflared` ships a diagnostic system collector for Linux, macOS, and
FreeBSD, but not OpenBSD, so recent releases fail to build for `openbsd/amd64`.
This repository supplies the missing collector and a script that cross-builds a
static binary from an unmodified upstream tree.

The tunnel data path is upstream code, untouched. The added collector runs only
under `cloudflared tunnel diag`.

## This repo or the OpenBSD port

Two projects add OpenBSD support to cloudflared using the same underlying patch.
Choose by how you want to install:

- [**openbsd-port-cloudflared**](https://github.com/ivoronin/openbsd-port-cloudflared)
  is a proper OpenBSD ports-tree port: `make install`, a packaged build, an
  `rc.d` service, and a dedicated user. If you are on OpenBSD and want a managed,
  packaged install done the standard way, use that.
- **This repo** is a standalone cross-build script. It produces a single static
  `openbsd/amd64` binary from any build host (Linux, macOS, a CI runner) with no
  ports tree required. Use it when you want a binary quickly, or to build OpenBSD
  binaries without an OpenBSD host.

This repo stays deliberately small and build-focused. For a maintained,
official-channel install, prefer the ports-tree port.

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

## Running as a service

The build produces a binary only. To run it as a managed daemon on OpenBSD:

1. Create an unprivileged user:

   ```sh
   doas useradd -d /var/empty -s /sbin/nologin -c "Cloudflare Tunnel" -L daemon _cloudflared
   ```

2. Configure the tunnel. cloudflared reads `/etc/cloudflared/config.yml`; point
   it at your tunnel and its credentials file (see Cloudflare's docs). Keep the
   token or credentials in that file, not in the rc script, and make it readable
   only by `_cloudflared`.

3. Install the rc.d script from [`contrib/cloudflared.rc`](contrib/cloudflared.rc):

   ```sh
   doas install -o root -g bin -m 555 contrib/cloudflared.rc /etc/rc.d/cloudflared
   doas rcctl enable cloudflared
   doas rcctl start cloudflared
   ```

For a named tunnel, override the flags in `/etc/rc.conf.local`, for example:

```sh
rcctl set cloudflared flags "--no-autoupdate tunnel run mytunnel"
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

This is a bridge until OpenBSD support lands in `cloudflared` itself. The same
patch is carried by the [ports-tree port](https://github.com/ivoronin/openbsd-port-cloudflared);
the durable home for OpenBSD support is upstream cloudflared or the ports tree,
not a standalone build script. If you can, help move it there.

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
