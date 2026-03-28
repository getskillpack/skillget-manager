package skillgetmanager

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPublishSkill_requiresToken(t *testing.T) {
	t.Setenv("SKILLGET_REGISTRY_TOKEN", "")
	t.Setenv("SKILLGET_TOKEN", "")
	err := PublishSkill(context.Background(), []byte(`{}`), []byte("x"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPublishSkill_success(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/skills" || r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") != "Bearer testtok" {
			http.Error(w, "auth", http.StatusUnauthorized)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.FormValue("manifest") == "" {
			http.Error(w, "manifest", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("archive")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		body, _ := io.ReadAll(file)
		if string(body) != "gz-bytes" {
			http.Error(w, "bad archive", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	t.Setenv("SKILLGET_REGISTRY_URL", ts.URL+"/api/v1")
	t.Setenv("SKILLGET_REGISTRY_TOKEN", "testtok")
	HTTPClient = ts.Client()

	err := PublishSkill(context.Background(), []byte(`{"name":"n","version":"1.0.0"}`), []byte("gz-bytes"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchSkills_authorQuery(t *testing.T) {
	prev := HTTPClient
	t.Cleanup(func() { HTTPClient = prev })

	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.RequestURI()
		_ = json.NewEncoder(w).Encode(ListSkillsResponse{})
	}))
	defer ts.Close()

	t.Setenv("SKILLGET_REGISTRY_URL", ts.URL+"/api/v1")
	HTTPClient = ts.Client()

	_, err := SearchSkills(context.Background(), SearchSkillsOptions{Author: "acme", Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotPath, "author=acme") {
		t.Fatalf("expected author in query, got %q", gotPath)
	}
}
