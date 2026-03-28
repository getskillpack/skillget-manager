package skillgetmanager

import (
	"context"
	"fmt"
	"net/url"
	"strings"
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

// pickInstallableVersion chooses a version for an unpinned install.
// Prefers registry latest_version when present and not marked yanked in versions[];
// otherwise the first non-yanked entry in versions (registry order).
func pickInstallableVersion(d SkillDetail) string {
	if d.LatestVersion != nil {
		candidate := strings.TrimSpace(*d.LatestVersion)
		if candidate != "" && !versionMarkedYanked(d, candidate) {
			return candidate
		}
	}
	for _, v := range d.Versions {
		if v.IsYanked {
			continue
		}
		if v.Version != "" {
			return v.Version
		}
	}
	return ""
}

// versionMarkedYanked reports whether v appears in d.Versions with is_yanked true.
// If v is not listed, returns false (trust latest_version from the registry).
func versionMarkedYanked(d SkillDetail, v string) bool {
	for _, row := range d.Versions {
		if row.Version == v {
			return row.IsYanked
		}
	}
	return false
}
