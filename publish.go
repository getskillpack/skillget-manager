package skillgetmanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

// PublishSkill uploads manifest JSON and a .tar.gz archive via multipart POST /skills.
// Requires a non-empty RegistryToken() (SKILLGET_REGISTRY_TOKEN or SKILLGET_TOKEN).
func PublishSkill(ctx context.Context, manifestJSON []byte, archive []byte) error {
	token := RegistryToken()
	if strings.TrimSpace(token) == "" {
		return fmt.Errorf("publish requires SKILLGET_REGISTRY_TOKEN (or SKILLGET_TOKEN)")
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	if err := mw.WriteField("manifest", string(manifestJSON)); err != nil {
		return err
	}
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", `form-data; name="archive"; filename="skill.tar.gz"`)
	hdr.Set("Content-Type", "application/gzip")
	part, err := mw.CreatePart(hdr)
	if err != nil {
		return err
	}
	if _, err := part.Write(archive); err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}

	url := joinURLPath(RegistryBaseURL(), "/skills")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("registry request failed: %w%s", err, registryNetworkHint())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusCreated {
		return nil
	}
	text := strings.TrimSpace(string(body))
	if text == "" {
		text = url
	}
	return fmt.Errorf("publish failed: registry %s: %s%s", res.Status, text, registryHintForStatus(res.StatusCode))
}
