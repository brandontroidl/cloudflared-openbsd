# Contributing

Thanks for your interest. This is a small, focused repository: its only job is to
keep `cloudflared` building and running on OpenBSD until upstream carries the
support itself. Contributions that keep the patch minimal and current are the
most useful.

## Scope

In scope:

- Fixing the OpenBSD system collector (`system_collector_openbsd.go`)
- Keeping `build.sh` working against new upstream cloudflared releases
- Documentation corrections

Out of scope:

- Changes to cloudflared's tunnel behavior, protocols, or dependencies. Those
  belong upstream at https://github.com/cloudflare/cloudflared.
- New platforms or architectures beyond `openbsd/amd64`, unless you can test
  them on real hardware.

The best outcome for any collector fix is that it also lands upstream, so the
patch here eventually disappears. If you can, open the equivalent change against
`cloudflare/cloudflared` too.

## Building and testing

```sh
sh build.sh            # default upstream version
sh build.sh <version>  # a specific tag
```

A change is not done until it is verified on a real OpenBSD host:

1. Build the binary (cross-building from any host is fine).
2. Copy it to an OpenBSD `amd64` machine and install it (see the README).
3. Run `cloudflared tunnel diag` and confirm the memory and file-descriptor
   fields populate without errors.

State which OpenBSD release and cloudflared version you tested against in your
pull request. "It compiles" is not sufficient; the collector must run.

## Style

- Match upstream cloudflared's Go conventions. Run `gofmt` on any Go file.
- Keep the `//go:build openbsd` constraint intact on the collector.
- Keep the patch surface small. Prefer adapting the existing FreeBSD approach
  over introducing new abstractions.
- No em-dashes in prose; use spaced hyphens or restructure.
- POSIX `sh` for scripts, not bash-isms.

## Submitting changes

1. Open an issue first for anything non-trivial, so the approach can be agreed
   before you invest time.
2. Keep each pull request to one logical change.
3. Describe what you changed, why, and how you tested it (release + version).

## Reporting bugs

Open an issue with the OpenBSD release, the cloudflared version, the exact
command, and the full output. For anything security-sensitive, follow
[`SECURITY.md`](SECURITY.md) instead of opening a public issue.

## License

By contributing, you agree that your contributions are licensed under Apache-2.0,
the same license as this repository and cloudflared.
