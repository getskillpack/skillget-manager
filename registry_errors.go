package skillgetmanager

func registryNetworkHint() string {
	return "\nHint: check network connectivity, DNS, and SKILLGET_REGISTRY_URL (path usually ends with /api/v1)."
}

func registryHintForStatus(status int) string {
	switch status {
	case 400:
		return "\nHint: request may be invalid — check query params and paths against the registry API."
	case 401:
		return "\nHint: set SKILLGET_REGISTRY_TOKEN (or SKILLGET_TOKEN) for authenticated requests."
	case 403:
		return "\nHint: token may lack permission, or the registry blocks this operation for anonymous clients."
	case 404:
		return "\nHint: check the skill name and SKILLGET_REGISTRY_URL."
	case 409:
		return "\nHint: this skill version already exists; bump version or yank the old release."
	case 410:
		return "\nHint: this version was yanked and cannot be installed."
	case 422:
		return "\nHint: manifest or form fields failed validation — compare with the registry API schema."
	case 429:
		return "\nHint: rate limited — back off and retry; reduce parallel requests if you batch installs or publishes."
	case 503:
		return "\nHint: registry write API may be disabled on the server."
	default:
		return ""
	}
}
