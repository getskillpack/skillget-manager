package skillgetmanager

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

// NameVersion holds a parsed skill spec "name" or "name@version".
type NameVersion struct {
	Name    string
	Version string // empty means "latest"
}

// ParseNameVersion splits spec at the last '@' (same rule as the TS prototype).
func ParseNameVersion(spec string) NameVersion {
	at := strings.LastIndex(spec, "@")
	if at <= 0 {
		return NameVersion{Name: spec}
	}
	return NameVersion{
		Name:    spec[:at],
		Version: spec[at+1:],
	}
}

// ResolveInstallTarget resolves a concrete version and archive metadata (no download).
func ResolveInstallTarget(ctx context.Context, spec string) (*VersionDetail, error) {
	nv := ParseNameVersion(spec)
	version := nv.Version
	if version == "" {
		var detail SkillDetail
		path := "/skills/" + url.PathEscape(nv.Name)
		if err := FetchJSON(ctx, path, &detail); err != nil {
			return nil, err
		}
		version = pickInstallableVersion(detail)
		if version == "" {
			return nil, fmt.Errorf("no installable versions for skill %q", nv.Name)
		}
	}
	var vd VersionDetail
	vpath := "/skills/" + url.PathEscape(nv.Name) + "/versions/" + url.PathEscape(version)
	if err := FetchJSON(ctx, vpath, &vd); err != nil {
		return nil, err
	}
	return &vd, nil
}

// pickInstallableVersion chooses the highest semver among non-yanked versions in the
// skill detail map, matching reference registry logic (see getskillpack/registry filestore).
func pickInstallableVersion(d SkillDetail) string {
	if len(d.Versions) == 0 {
		return ""
	}
	var keys []string
	for v, info := range d.Versions {
		v = strings.TrimSpace(v)
		if v == "" || info.Yanked {
			continue
		}
		keys = append(keys, v)
	}
	if len(keys) == 0 {
		return ""
	}
	sort.Slice(keys, func(i, j int) bool {
		return semver.Compare(canonicalSemver(keys[i]), canonicalSemver(keys[j])) < 0
	})
	return keys[len(keys)-1]
}

func canonicalSemver(v string) string {
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}
