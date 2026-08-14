package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadServiceManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte("entries:\n  - one.service\n  - two.service\n"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	manifest, err := LoadServiceManifest(path)
	if err != nil {
		t.Fatalf("LoadServiceManifest failed: %v", err)
	}
	if len(manifest.Entries) != 2 || manifest.Entries[0] != "one.service" || manifest.Entries[1] != "two.service" {
		t.Fatalf("unexpected entries: %#v", manifest.Entries)
	}
}
