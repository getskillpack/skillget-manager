package skillgetmanager

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var checksumSHA256 = regexp.MustCompile(`^sha256:([a-fA-F0-9]{64})$`)

// DownloadSkillResult is returned after a successful archive download.
type DownloadSkillResult struct {
	ArchivePath string
	Meta        struct {
		Name     string
		Version  string
		Checksum string
	}
}

// DownloadSkillOptions configures download paths.
type DownloadSkillOptions struct {
	Cwd        string // default: os.Getwd()
	OutputPath string // if empty: .skillget/skills/<name>/<version>/<name>-<version>.tar.gz under Cwd
}

// DownloadSkillArchive downloads the skill archive and updates skills.lock in cwd.
func DownloadSkillArchive(ctx context.Context, spec string, opts DownloadSkillOptions) (*DownloadSkillResult, error) {
	cwd := opts.Cwd
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	meta, err := ResolveInstallTarget(ctx, spec)
	if err != nil {
		return nil, err
	}

	fileName := meta.Name + "-" + meta.Version + ".tar.gz"
	dest := opts.OutputPath
	if dest == "" {
		dest = filepath.Join(cwd, ".skillget", "skills", meta.Name, meta.Version, fileName)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, meta.ArchiveURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("download failed %s: %s", res.Status, meta.ArchiveURL)
	}

	f, err := os.Create(dest)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(f, res.Body); err != nil {
		_ = os.Remove(dest)
		return nil, err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(dest)
		return nil, err
	}

	if err := verifyDownloadSHA256(dest, meta.Checksum); err != nil {
		_ = os.Remove(dest)
		return nil, err
	}

	lock := ReadSkillsLock(cwd)
	if lock.Skills == nil {
		lock.Skills = map[string]string{}
	}
	lock.Skills[meta.Name] = meta.Version
	if err := WriteSkillsLock(cwd, lock); err != nil {
		return nil, err
	}

	out := &DownloadSkillResult{ArchivePath: dest}
	out.Meta.Name = meta.Name
	out.Meta.Version = meta.Version
	out.Meta.Checksum = meta.Checksum
	return out, nil
}

func verifyDownloadSHA256(path, want string) error {
	want = strings.TrimSpace(want)
	if want == "" {
		return nil
	}
	m := checksumSHA256.FindStringSubmatch(want)
	if m == nil {
		return fmt.Errorf("unsupported checksum format %q (expected sha256:<64 hex>)", want)
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(got, m[1]) {
		return fmt.Errorf("checksum mismatch for %s: expected %s, got %s", filepath.Base(path), m[1], got)
	}
	return nil
}
