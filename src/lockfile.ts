import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import type { SkillsLockfile } from "./types.js";

export const LOCKFILE_NAME = "skills.lock";

export function readSkillsLock(cwd: string): SkillsLockfile {
  const p = join(cwd, LOCKFILE_NAME);
  try {
    const raw = readFileSync(p, "utf8");
    const j = JSON.parse(raw) as Partial<SkillsLockfile>;
    if (j?.lockfileVersion === 1 && j.skills && typeof j.skills === "object" && !Array.isArray(j.skills)) {
      return { lockfileVersion: 1, skills: { ...j.skills } };
    }
  } catch {
    /* missing or invalid */
  }
  return { lockfileVersion: 1, skills: {} };
}

export function writeSkillsLock(cwd: string, lock: SkillsLockfile): void {
  const p = join(cwd, LOCKFILE_NAME);
  writeFileSync(p, `${JSON.stringify(lock, null, 2)}\n`, "utf8");
}
