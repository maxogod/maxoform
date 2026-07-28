package dconf

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/logger"
	"go.uber.org/zap"
)

type mockRunner struct {
	name      string
	args      []string
	stdinData string
	runErr    error
}

func (m *mockRunner) Run(_ string, _ ...string) error { return nil }
func (m *mockRunner) RunShell(_ string) error         { return nil }
func (m *mockRunner) Check(_ string, _ ...string) bool {
	return false
}
func (m *mockRunner) CheckShell(_ string) bool { return false }
func (m *mockRunner) RunWithStdin(stdin io.Reader, name string, args ...string) error {
	data, _ := io.ReadAll(stdin)
	m.stdinData = string(data)
	m.name = name
	m.args = args
	return m.runErr
}

func TestImport_LoadsAllManifestEntries(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "a.ini"), []byte("[x]\nfoo='bar'\n"), 0o644); err != nil {
		t.Fatalf("write settings file: %v", err)
	}

	manifest := &config.DconfManifest{
		Entries: []config.DconfEntry{
			{Key: "/org/gnome/test/", File: "a.ini"},
		},
	}
	m := &mockRunner{}
	if err := Import(m, dir, manifest); err != nil {
		t.Fatalf("Import failed: %v", err)
	}
	if m.name != "dconf" || len(m.args) != 2 || m.args[0] != "load" || m.args[1] != "/org/gnome/test/" {
		t.Fatalf("unexpected command: %s %#v", m.name, m.args)
	}
	if m.stdinData == "" {
		t.Fatalf("expected stdin content to be forwarded")
	}
}

func TestImport_ReturnsOpenError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	manifest := &config.DconfManifest{
		Entries: []config.DconfEntry{
			{Key: "/org/gnome/test/", File: "missing.ini"},
		},
	}
	if err := Import(&mockRunner{}, t.TempDir(), manifest); err == nil {
		t.Fatalf("expected error for missing file")
	}
}
