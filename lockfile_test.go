package skillgetmanager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadWriteSkillsLock(t *testing.T) {
	dir := t.TempDir()
	lock := ReadSkillsLock(dir)
	if lock.LockfileVersion != 1 || len(lock.Skills) != 0 {
		t.Fatalf("empty read: %+v", lock)
	}
	lock.Skills["a"] = "1.0.0"
	if err := WriteSkillsLock(dir, lock); err != nil {
		t.Fatal(err)
	}
	again := ReadSkillsLock(dir)
	if again.Skills["a"] != "1.0.0" {
		t.Fatalf("roundtrip: %+v", again)
	}
}

func TestReadSkillsLockInvalidSkillsType(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, LockfileName)
	if err := os.WriteFile(p, []byte(`{"lockfileVersion":1,"skills":{"x":1}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	lock := ReadSkillsLock(dir)
	if len(lock.Skills) != 0 {
		t.Fatalf("expected empty skills on invalid types, got %+v", lock)
	}
}
