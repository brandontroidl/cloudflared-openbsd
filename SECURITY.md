
  # Security Policy

  ## Scope

  This repository contains only:

  - `system_collector_openbsd.go`, an OpenBSD implementation of cloudflared's
    diagnostic `SystemCollector` (sysctl reads, no privileged operations)
  - `build.sh`, which clones upstream cloudflared at a tagged version, applies
    the patch, and cross-builds a static `openbsd/amd64` binary

  Only issues in **these files** are in scope here.

  **Vulnerabilities in cloudflared itself** (the tunnel, its protocols, its
  dependencies) are out of scope. Report those to Cloudflare:
  https://github.com/cloudflare/cloudflared/security or
  https://hackerone.com/cloudflare

  This is not an official Cloudflare project.

  ## Supported versions

  Only the latest commit on the default branch is supported. The build script
  takes an upstream version argument; binaries built from cloudflared versions
  that Cloudflare no longer supports inherit upstream's risk regardless of
  this patch.

  ## Reporting a vulnerability

  Use GitHub's private vulnerability reporting on this repository
  (Security tab -> "Report a vulnerability"). Please do not open a public
  issue for security reports.

  Include the affected file, the impact, and reproduction steps. You can
  expect an initial response within 7 days.

  ## Threat model notes for users

  - `build.sh` fetches upstream source over HTTPS from
    `github.com/cloudflare/cloudflared` at the tag you specify. Review the
    script before running it; it builds with `CGO_ENABLED=0` and needs no
    root privileges.
  - The patched code runs only when you invoke `cloudflared tunnel diag`.
    The tunnel data path is unmodified upstream code.
  - Install the resulting binary with the ownership and mode shown in the
    README (`root:bin`, `0755`); the binary itself needs no special
    privileges beyond what cloudflared normally requires.
