# Contributing to skillget-manager

Library used by the official CLI ([`getskillpack/cli`](https://github.com/getskillpack/cli)). End-user install flows and tutorials live in that repo (`docs/`).

## Prerequisites

- **Go 1.22+**
- A running **registry** when exercising integration-style tests (see below).

## Clone and test

```bash
git clone https://github.com/getskillpack/skillget-manager.git
cd skillget-manager
go test ./...
```

## Integration tests

Some tests talk to a real registry when `SKILLGET_REGISTRY_URL` points at a reachable API (see `resolve_integration_test.go`). To skip them locally:

```bash
go test ./... -short
```

CI should run the full suite when the registry URL and credentials are available.

## Work alongside the CLI

During active development, `getskillpack/cli` often uses a `replace` in `go.mod` pointing at a local clone of this module. After tagged releases are visible to the Go module proxy, trim or remove `replace` and bump the required version in `cli`.

Registry contract: [`getskillpack/registry` API](https://github.com/getskillpack/registry/blob/main/docs/registry-api.md).

## Releases

- Follow [CHANGELOG.md](CHANGELOG.md) and semver for public tags.
- Keep [CHANGELOG.md](CHANGELOG.md) **Unreleased** updated as you merge user-visible behavior changes.

## Security

See [SECURITY.md](SECURITY.md) for coordinated disclosure.
