package skillgetmanager

// SkillsLockfile is the on-disk lockfile shape (v1).
type SkillsLockfile struct {
	LockfileVersion int               `json:"lockfileVersion"`
	Skills          map[string]string `json:"skills"`
}

// ListSkillsResponse is the registry list/search JSON envelope.
type ListSkillsResponse struct {
	Data []struct {
		Name           string `json:"name"`
		Description    string `json:"description,omitempty"`
		Author         string `json:"author,omitempty"`
		LatestVersion  string `json:"latest_version,omitempty"`
	} `json:"data"`
	Meta *struct {
		Total int `json:"total,omitempty"`
	} `json:"meta,omitempty"`
}

// VersionDetail is returned for a concrete skill version (install target).
type VersionDetail struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	ArchiveURL string `json:"archive_url"`
	Checksum   string `json:"checksum,omitempty"`
}

// SkillDetail is skill metadata from the registry.
type SkillDetail struct {
	Name           string `json:"name"`
	RepositoryURL  *string `json:"repository_url"`
	Homepage       *string `json:"homepage"`
	Dependencies   []struct {
		Name  string `json:"name"`
		Range string `json:"range,omitempty"`
	} `json:"dependencies,omitempty"`
	LatestVersion *string `json:"latest_version"`
	Versions      []struct {
		Version   string `json:"version"`
		IsYanked  bool   `json:"is_yanked,omitempty"`
	} `json:"versions,omitempty"`
}
