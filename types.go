package skillgetmanager

import (
	"encoding/json"
	"time"
)

// SkillsLockfile is the on-disk lockfile shape (v1).
type SkillsLockfile struct {
	LockfileVersion int               `json:"lockfileVersion"`
	Skills          map[string]string `json:"skills"`
}

// ListSkillsResponse is the registry list/search JSON envelope.
type ListSkillsResponse struct {
	Data []struct {
		Name          string    `json:"name"`
		Description   string    `json:"description,omitempty"`
		Author        string    `json:"author,omitempty"`
		LatestVersion string    `json:"latest_version,omitempty"`
		CreatedAt     time.Time `json:"created_at,omitempty"`
	} `json:"data"`
	Meta *struct {
		Total  int `json:"total,omitempty"`
		Limit  int `json:"limit,omitempty"`
		Offset int `json:"offset,omitempty"`
	} `json:"meta,omitempty"`
}

// VersionPublicInfo matches GET /skills/{name} → versions[version] in registry-api.md.
type VersionPublicInfo struct {
	Manifest    json.RawMessage `json:"manifest,omitempty"`
	Checksum    string          `json:"checksum,omitempty"`
	ArchiveURL  string          `json:"archive_url,omitempty"`
	PublishedAt time.Time       `json:"published_at,omitempty"`
	Yanked      bool            `json:"yanked,omitempty"`
}

// SkillDetail is returned by GET /skills/{name} (registry compiled-core contract).
type SkillDetail struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description,omitempty"`
	Author      string                       `json:"author,omitempty"`
	CreatedAt   time.Time                    `json:"created_at,omitempty"`
	Versions    map[string]VersionPublicInfo `json:"versions,omitempty"`
}

// VersionDetail is returned for GET /skills/{name}/versions/{version} (install target).
type VersionDetail struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	ArchiveURL string `json:"archive_url"`
	Checksum   string `json:"checksum,omitempty"`
}
