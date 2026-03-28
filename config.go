package skillgetmanager

import (
	"os"
	"strings"
)

const defaultRegistryBase = "https://registry.skpkg.org/api/v1"

func trimBase(raw string) string {
	t := strings.TrimSpace(raw)
	if t == "" {
		return ""
	}
	return strings.TrimSuffix(t, "/")
}

// RegistryBaseURL returns the registry API base URL.
// SKILLGET_REGISTRY_URL wins; SKPKG_REGISTRY_URL is a legacy fallback.
func RegistryBaseURL() string {
	if u := trimBase(os.Getenv("SKILLGET_REGISTRY_URL")); u != "" {
		return u
	}
	if u := trimBase(os.Getenv("SKPKG_REGISTRY_URL")); u != "" {
		return u
	}
	return defaultRegistryBase
}

// RegistryConfigSource reports which setting selected the registry URL.
func RegistryConfigSource() string {
	if trimBase(os.Getenv("SKILLGET_REGISTRY_URL")) != "" {
		return "SKILLGET_REGISTRY_URL"
	}
	if trimBase(os.Getenv("SKPKG_REGISTRY_URL")) != "" {
		return "SKPKG_REGISTRY_URL"
	}
	return "default"
}

func trimToken(raw string) string {
	t := strings.TrimSpace(raw)
	if t == "" {
		return ""
	}
	return t
}

// RegistryToken returns the bearer token for authenticated registry calls (e.g. publish).
// SKILLGET_REGISTRY_TOKEN wins; SKILLGET_TOKEN is a legacy alias.
func RegistryToken() string {
	if t := trimToken(os.Getenv("SKILLGET_REGISTRY_TOKEN")); t != "" {
		return t
	}
	return trimToken(os.Getenv("SKILLGET_TOKEN"))
}
