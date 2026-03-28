import { fetchJson } from "./client.js";
import type { ListSkillsResponse } from "./types.js";

export type SearchSkillsOptions = {
  query?: string;
  limit?: number;
  offset?: number;
};

export async function searchSkills(opts: SearchSkillsOptions = {}): Promise<ListSkillsResponse> {
  const limit = opts.limit ?? 20;
  const offset = opts.offset ?? 0;
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  if (opts.query) params.set("q", opts.query);
  return fetchJson<ListSkillsResponse>(`/skills?${params.toString()}`);
}
