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

`getskillpack/cli` pins this module by **semver** in `go.mod`. For a local loop you can still add a temporary `replace` → sibling checkout. With **private** org repos, CLI contributors and CI use `GOPRIVATE` / `GONOSUMDB` and Git credentials (see `cli` repo `docs/BOARD_PAT_QUICK_RU.md` § 4). When the module is public on the Go proxy, drop `GOPRIVATE` for this path and use plain `go get …@vX.Y.Z`.

Registry contract: [`getskillpack/registry` API](https://github.com/getskillpack/registry/blob/main/docs/registry-api.md).

## Releases

- Follow [CHANGELOG.md](CHANGELOG.md) and semver for public tags.
- Keep [CHANGELOG.md](CHANGELOG.md) **Unreleased** updated as you merge user-visible behavior changes.

## Security

See [SECURITY.md](SECURITY.md) for coordinated disclosure.
