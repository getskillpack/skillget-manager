import { fetchJson } from "./client.js";
import type { SkillDetail, VersionDetail } from "./types.js";

export function parseNameVersion(spec: string): { name: string; version?: string } {
  const at = spec.lastIndexOf("@");
  if (at <= 0) return { name: spec };
  return { name: spec.slice(0, at), version: spec.slice(at + 1) };
}

/** Resolve concrete version and fetch archive metadata (does not download bytes). */
export async function resolveInstallTarget(spec: string): Promise<VersionDetail> {
  const { name, version: pinned } = parseNameVersion(spec);
  let version = pinned;
  if (!version) {
    const detail = await fetchJson<SkillDetail>(`/skills/${encodeURIComponent(name)}`);
    const active = detail.versions?.filter((v) => !v.is_yanked) ?? [];
    const latest = active[0]?.version;
    if (!latest) {
      throw new Error(`No installable versions for skill "${name}".`);
    }
    version = latest;
  }
  return fetchJson<VersionDetail>(
    `/skills/${encodeURIComponent(name)}/versions/${encodeURIComponent(version)}`,
  );
}
