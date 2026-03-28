package skillgetmanager

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveInstallTarget_latest(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/skills/demo":
			_ = json.NewEncoder(w).Encode(SkillDetail{
				Name: "demo",
				Versions: []struct {
					Version  string `json:"version"`
					IsYanked bool   `json:"is_yanked,omitempty"`
				}{
					{Version: "2.0.0", IsYanked: true},
					{Version: "1.0.0"},
				},
			})
		case "/api/v1/skills/demo/versions/1.0.0":
			_ = json.NewEncoder(w).Encode(VersionDetail{
				Name:       "demo",
				Version:    "1.0.0",
				ArchiveURL: "https://example.invalid/a.tgz",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	t.Setenv("SKILLGET_REGISTRY_URL", ts.URL+"/api/v1")
	HTTPClient = ts.Client()

	vd, err := ResolveInstallTarget(context.Background(), "demo")
	if err != nil {
		t.Fatal(err)
	}
	if vd.Version != "1.0.0" || vd.ArchiveURL == "" {
		t.Fatalf("unexpected: %+v", vd)
	}
}

func TestResolveInstallTarget_pinned(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/skills/demo/versions/3.4.5" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(VersionDetail{
			Name:       "demo",
			Version:    "3.4.5",
			ArchiveURL: "https://example.invalid/b.tgz",
		})
	}))
	defer ts.Close()

	t.Setenv("SKILLGET_REGISTRY_URL", ts.URL+"/api/v1")
	HTTPClient = ts.Client()

	vd, err := ResolveInstallTarget(context.Background(), "demo@3.4.5")
	if err != nil {
		t.Fatal(err)
	}
	if vd.Version != "3.4.5" {
		t.Fatalf("unexpected: %+v", vd)
	}
}

func TestResolveInstallTarget_latestVersionField(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/skills/demo":
			lv := "9.9.9"
			_ = json.NewEncoder(w).Encode(SkillDetail{
				Name:          "demo",
				LatestVersion: &lv,
				Versions: []struct {
					Version  string `json:"version"`
					IsYanked bool   `json:"is_yanked,omitempty"`
				}{
					{Version: "0.0.1"},
				},
			})
		case "/api/v1/skills/demo/versions/9.9.9":
			_ = json.NewEncoder(w).Encode(VersionDetail{
				Name:       "demo",
				Version:    "9.9.9",
				ArchiveURL: "https://example.invalid/c.tgz",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	t.Setenv("SKILLGET_REGISTRY_URL", ts.URL+"/api/v1")
	HTTPClient = ts.Client()

	vd, err := ResolveInstallTarget(context.Background(), "demo")
	if err != nil {
		t.Fatal(err)
	}
	if vd.Version != "9.9.9" {
		t.Fatalf("expected latest_version 9.9.9, got %+v", vd)
	}
}
