export type SkillsLockfile = {
  lockfileVersion: 1;
  skills: Record<string, string>;
};

export type ListSkillsResponse = {
  data: Array<{
    name: string;
    description?: string;
    author?: string;
    latest_version?: string;
  }>;
  meta?: { total?: number };
};

export type VersionDetail = {
  name: string;
  version: string;
  archive_url: string;
  checksum?: string;
};

export type SkillDetail = {
  name: string;
  repository_url?: string | null;
  homepage?: string | null;
  dependencies?: Array<{ name: string; range?: string }>;
  latest_version?: string | null;
  versions?: Array<{ version: string; is_yanked?: boolean }>;
};
