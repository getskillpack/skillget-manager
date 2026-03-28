package skillgetmanager

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadSkillArchive_checksumOK(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	payload := []byte("fake-gzip-payload")
	sum := sha256.Sum256(payload)
	checksum := "sha256:" + hex.EncodeToString(sum[:])

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/skills/demo":
			lv := "1.0.0"
			_ = json.NewEncoder(w).Encode(SkillDetail{Name: "demo", LatestVersion: &lv})
		case "/api/v1/skills/demo/versions/1.0.0":
			_ = json.NewEncoder(w).Encode(VersionDetail{
				Name:       "demo",
				Version:    "1.0.0",
				ArchiveURL: "https://" + r.Host + "/blob.tgz",
				Checksum:   checksum,
			})
		case "/blob.tgz":
			_, _ = w.Write(payload)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	t.Setenv("SKILLGET_REGISTRY_URL", ts.URL+"/api/v1")
	HTTPClient = ts.Client()

	dir := t.TempDir()
	res, err := DownloadSkillArchive(context.Background(), "demo", DownloadSkillOptions{Cwd: dir})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(res.ArchivePath); err != nil {
		t.Fatal(err)
	}
}

func TestDownloadSkillArchive_checksumMismatch(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/skills/demo":
			lv := "1.0.0"
			_ = json.NewEncoder(w).Encode(SkillDetail{Name: "demo", LatestVersion: &lv})
		case "/api/v1/skills/demo/versions/1.0.0":
			_ = json.NewEncoder(w).Encode(VersionDetail{
				Name:       "demo",
				Version:    "1.0.0",
				ArchiveURL: "https://" + r.Host + "/blob.tgz",
				Checksum:   "sha256:0000000000000000000000000000000000000000000000000000000000000000",
			})
		case "/blob.tgz":
			_, _ = w.Write([]byte("not-zero-hash"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	t.Setenv("SKILLGET_REGISTRY_URL", ts.URL+"/api/v1")
	HTTPClient = ts.Client()

	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.tgz")
	_, err := DownloadSkillArchive(context.Background(), "demo", DownloadSkillOptions{Cwd: dir, OutputPath: outPath})
	if err == nil {
		t.Fatal("expected checksum error")
	}
	if _, err := os.Stat(outPath); err == nil {
		t.Fatal("expected archive removed on checksum failure")
	}
}
