package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadServiceManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte("entries:\n  - condition: test -f /tmp/a\n    service: one.service\n  - condition: test -f /tmp/b\n    service: two.service\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	manifest, err := LoadServiceManifest(path)
	if err != nil {
		t.Fatalf("LoadServiceManifest failed: %v", err)
	}
	if len(manifest.Entries) != 2 ||
		manifest.Entries[0].Condition != "test -f /tmp/a" ||
		manifest.Entries[0].Service != "one.service" ||
		manifest.Entries[1].Condition != "test -f /tmp/b" ||
		manifest.Entries[1].Service != "two.service" {
		t.Fatalf("unexpected entries: %#v", manifest.Entries)
	}
}
