# Security Policy

## Scope

This repository contains only:

- `system_collector_openbsd.go` - an OpenBSD implementation of cloudflared's
  diagnostic `SystemCollector`. It runs read-only `sysctl` commands and performs
  no privileged operations.
- `build.sh` - clones upstream cloudflared at a tagged version, applies the
  patch, and cross-builds a static `openbsd/amd64` binary.

Only issues in **these files** are in scope for this repository.

Vulnerabilities in cloudflared itself (the tunnel, its protocols, its
dependencies) are **out of scope here**. Report those directly to Cloudflare:

- https://github.com/cloudflare/cloudflared/security
- https://hackerone.com/cloudflare

This is not an official Cloudflare project.

## Supported versions

Only the latest commit on the default branch is supported. `build.sh` accepts an
upstream version argument; a binary built from a cloudflared version that
Cloudflare no longer supports inherits upstream's risk regardless of this patch.

## Reporting a vulnerability

Use GitHub's private vulnerability reporting on this repository: open the
**Security** tab and choose **Report a vulnerability**. Do not open a public
issue for a security report.

Please include the affected file, the impact, and steps to reproduce. Expect an
initial response within 7 days.

## Threat model notes for users

- `build.sh` fetches upstream source over HTTPS from
  `github.com/cloudflare/cloudflared` at the tag you specify. Review the script
  before running it. It builds with `CGO_ENABLED=0` and needs no root
  privileges.
- The patched code runs only when you invoke `cloudflared tunnel diag`. The
  tunnel data path is unmodified upstream code.
- Install the resulting binary with the ownership and mode shown in the README
  (`root:bin`, `0755`). The binary needs no privileges beyond what cloudflared
  normally requires.
