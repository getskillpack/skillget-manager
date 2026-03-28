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
		if detail.LatestVersion != nil && *detail.LatestVersion != "" {
			version = *detail.LatestVersion
		} else {
			var latest string
			for _, v := range detail.Versions {
				if v.IsYanked {
					continue
				}
				latest = v.Version
				break
			}
			if latest == "" {
				return nil, fmt.Errorf("no installable versions for skill %q", nv.Name)
			}
			version = latest
		}
	}
	var vd VersionDetail
	vpath := "/skills/" + url.PathEscape(nv.Name) + "/versions/" + url.PathEscape(version)
	if err := FetchJSON(ctx, vpath, &vd); err != nil {
		return nil, err
	}
	return &vd, nil
}
