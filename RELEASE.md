# Releasing skillget-manager

Pre-1.0 library: tag **semver** on `main` when the API or behavior is stable enough for CLI and integrators to pin.

## Steps

1. Update [CHANGELOG.md](CHANGELOG.md): move items from **Unreleased** into a dated section `## 0.x.y — YYYY-MM-DD`.
2. Commit the changelog (and any release notes tweaks) on `main`.
3. Tag an **annotated** tag from that commit:

   ```bash
   git tag -a v0.1.0 -m "skillget-manager v0.1.0"
   git push origin main
   git push origin v0.1.0
   ```

4. Downstream ([getskillpack/cli](https://github.com/getskillpack/cli)): run `go get github.com/getskillpack/skillget-manager@v0.1.0`, refresh `go.sum`, remove the temporary `replace` stanza in `go.mod` once `proxy.golang.org` (or your `GOPRIVATE` + sum setup) resolves the module.

5. Align with [registry API](https://github.com/getskillpack/registry/blob/main/API.md) and CLI [RELEASE](https://github.com/getskillpack/cli/blob/main/docs/RELEASE.md) if HTTP contracts or install behavior changed.

CI on `main` must be green (`go test ./...`) before tagging.
