const DEFAULT_BASE = "https://registry.skpkg.org/api/v1";

function trimBase(raw: string | undefined): string | null {
  const t = raw?.trim();
  return t && t.length > 0 ? t.replace(/\/$/, "") : null;
}

/** Registry API base URL. `SKILLGET_REGISTRY_URL` wins; `SKPKG_REGISTRY_URL` is legacy fallback. */
export function registryBaseUrl(): string {
  return (
    trimBase(process.env.SKILLGET_REGISTRY_URL) ??
    trimBase(process.env.SKPKG_REGISTRY_URL) ??
    DEFAULT_BASE
  );
}

/** Which env var (if any) selected the registry URL. */
export function registryConfigSource(): "SKILLGET_REGISTRY_URL" | "SKPKG_REGISTRY_URL" | "default" {
  if (trimBase(process.env.SKILLGET_REGISTRY_URL)) return "SKILLGET_REGISTRY_URL";
  if (trimBase(process.env.SKPKG_REGISTRY_URL)) return "SKPKG_REGISTRY_URL";
  return "default";
}
