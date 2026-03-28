import { createWriteStream } from "node:fs";
import { mkdir } from "node:fs/promises";
import { dirname, join } from "node:path";
import { pipeline } from "node:stream/promises";
import { readSkillsLock, writeSkillsLock } from "./lockfile.js";
import { resolveInstallTarget } from "./resolve.js";

export type DownloadSkillResult = {
  archivePath: string;
  meta: { name: string; version: string; checksum?: string };
};

/**
 * Download skill archive, update skills.lock in cwd.
 * Default output: ./.skillget/skills/<name>/<version>/<name>-<version>.tar.gz
 */
export async function downloadSkillArchive(
  spec: string,
  options: { cwd?: string; outputPath?: string } = {},
): Promise<DownloadSkillResult> {
  const cwd = options.cwd ?? process.cwd();
  const meta = await resolveInstallTarget(spec);
  const fileName = `${meta.name}-${meta.version}.tar.gz`;
  const dest =
    options.outputPath ??
    join(cwd, ".skillget", "skills", meta.name, meta.version, fileName);
  await mkdir(dirname(dest), { recursive: true });

  const archiveRes = await fetch(meta.archive_url);
  if (!archiveRes.ok) {
    throw new Error(`Download failed ${archiveRes.status}: ${meta.archive_url}`);
  }
  if (!archiveRes.body) throw new Error("Empty response body");
  await pipeline(archiveRes.body, createWriteStream(dest));

  const lock = readSkillsLock(cwd);
  lock.skills[meta.name] = meta.version;
  writeSkillsLock(cwd, lock);

  return {
    archivePath: dest,
    meta: { name: meta.name, version: meta.version, checksum: meta.checksum },
  };
}
