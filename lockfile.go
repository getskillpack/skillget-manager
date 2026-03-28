package skillgetmanager

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LockfileName is the default lockfile filename in cwd.
const LockfileName = "skills.lock"

// ReadSkillsLock reads skills.lock from cwd; missing or invalid returns an empty v1 lock.
func ReadSkillsLock(cwd string) SkillsLockfile {
	p := filepath.Join(cwd, LockfileName)
	raw, err := os.ReadFile(p)
	if err != nil {
		return SkillsLockfile{LockfileVersion: 1, Skills: map[string]string{}}
	}
	var partial struct {
		LockfileVersion int                    `json:"lockfileVersion"`
		Skills          map[string]interface{} `json:"skills"`
	}
	if err := json.Unmarshal(raw, &partial); err != nil {
		return SkillsLockfile{LockfileVersion: 1, Skills: map[string]string{}}
	}
	if partial.LockfileVersion != 1 || partial.Skills == nil {
		return SkillsLockfile{LockfileVersion: 1, Skills: map[string]string{}}
	}
	out := make(map[string]string, len(partial.Skills))
	for k, v := range partial.Skills {
		s, ok := v.(string)
		if !ok {
			return SkillsLockfile{LockfileVersion: 1, Skills: map[string]string{}}
		}
		out[k] = s
	}
	return SkillsLockfile{LockfileVersion: 1, Skills: out}
}

// WriteSkillsLock writes skills.lock to cwd.
func WriteSkillsLock(cwd string, lock SkillsLockfile) error {
	p := filepath.Join(cwd, LockfileName)
	if lock.Skills == nil {
		lock.Skills = map[string]string{}
	}
	lock.LockfileVersion = 1
	raw, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(p, raw, 0o644)
}
