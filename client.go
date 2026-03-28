package skillgetmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTPClient is used for registry and archive requests. Tests may replace it.
var HTTPClient = http.DefaultClient

func joinURLPath(base, path string) string {
	p := path
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return base + p
}

// FetchJSON performs a GET for JSON from the registry API path (relative to RegistryBaseURL).
func FetchJSON(ctx context.Context, path string, out any) error {
	url := joinURLPath(RegistryBaseURL(), path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	res, err := HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		text := strings.TrimSpace(string(body))
		if text == "" {
			text = url
		}
		return fmt.Errorf("registry %s %s: %s", res.Status, res.Status, text)
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}
