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
	"strings"
	"testing"
	"time"
)

func TestDownloadSkillArchive_checksumOK(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	payload := []byte("fake-gzip-payload")
	sum := sha256.Sum256(payload)
	checksum := "sha256:" + hex.EncodeToString(sum[:])

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/skills/demo":
			_ = json.NewEncoder(w).Encode(SkillDetail{
				Name: "demo",
				Versions: map[string]VersionPublicInfo{
					"1.0.0": {Yanked: false},
				},
			})
		case "/api/v1/skills/demo/versions/1.0.0":
			_ = json.NewEncoder(w).Encode(VersionDetail{
				Name:       "demo",
				Version:    "1.0.0",
				ArchiveURL: ts.URL + "/blob.tgz",
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

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/skills/demo":
			_ = json.NewEncoder(w).Encode(SkillDetail{
				Name: "demo",
				Versions: map[string]VersionPublicInfo{
					"1.0.0": {},
				},
			})
		case "/api/v1/skills/demo/versions/1.0.0":
			_ = json.NewEncoder(w).Encode(VersionDetail{
				Name:       "demo",
				Version:    "1.0.0",
				ArchiveURL: ts.URL + "/blob.tgz",
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

func TestDownloadSkillArchive_sendsReadBearerToRegistry(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	var sawAuth string
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/skills/demo":
			sawAuth = r.Header.Get("Authorization")
			_ = json.NewEncoder(w).Encode(SkillDetail{
				Name: "demo",
				Versions: map[string]VersionPublicInfo{"1.0.0": {}},
			})
		case "/api/v1/skills/demo/versions/1.0.0":
			sawAuth = r.Header.Get("Authorization")
			_ = json.NewEncoder(w).Encode(VersionDetail{
				Name:       "demo",
				Version:    "1.0.0",
				ArchiveURL: ts.URL + "/blob.tgz",
			})
		case "/blob.tgz":
			if r.Header.Get("Authorization") != "Bearer readtok" {
				http.Error(w, "no auth", http.StatusUnauthorized)
				return
			}
			_, _ = w.Write([]byte("x"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	t.Setenv("SKILLGET_REGISTRY_URL", ts.URL+"/api/v1")
	t.Setenv("SKILLGET_REGISTRY_READ_TOKEN", "readtok")
	HTTPClient = ts.Client()

	dir := t.TempDir()
	_, err := DownloadSkillArchive(context.Background(), "demo", DownloadSkillOptions{Cwd: dir})
	if err != nil {
		t.Fatal(err)
	}
	if sawAuth != "Bearer readtok" {
		t.Fatalf("expected Bearer on registry requests, got %q", sawAuth)
	}
}

func TestSkillDetail_decodesRegistryShape(t *testing.T) {
	raw := `{
		"name": "x",
		"description": "d",
		"author": "a",
		"created_at": "2026-03-27T00:00:00Z",
		"versions": {
			"1.0.0": {
				"manifest": {"name":"x","version":"1.0.0"},
				"checksum": "sha256:abababababababababababababababababababababababababababababababab",
				"archive_url": "https://registry.example/downloads/ab.tar.gz",
				"published_at": "2026-03-27T00:00:00Z",
				"yanked": false
			}
		}
	}`
	var d SkillDetail
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		t.Fatal(err)
	}
	if d.Name != "x" || d.Description != "d" || d.Author != "a" {
		t.Fatalf("top-level: %+v", d)
	}
	if !d.CreatedAt.Equal(time.Date(2026, 3, 27, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("created_at: %v", d.CreatedAt)
	}
	v, ok := d.Versions["1.0.0"]
	if !ok {
		t.Fatal("missing version")
	}
	if v.Yanked || !strings.HasPrefix(string(v.Manifest), "{") {
		t.Fatalf("version info: %+v", v)
	}
}
