package services

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
	"go.uber.org/zap"
)

func TestInstallAndEnable_CopiesAndEnablesServices(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()

	dir := t.TempDir()
	serviceFile := "custom.service"
	if err := os.WriteFile(filepath.Join(dir, serviceFile), []byte("[Unit]\nDescription=Custom\n"), 0o644); err != nil {
		t.Fatalf("write service file: %v", err)
	}

	m := &shell.MockExecutor{}
	manifest := &config.ServiceManifest{
		Entries: []string{serviceFile},
	}

	if err := InstallAndEnable(m, dir, manifest); err != nil {
		t.Fatalf("InstallAndEnable failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{
		"sudo cp " + filepath.Join(dir, serviceFile) + " /etc/systemd/system/" + serviceFile,
		"sudo systemctl daemon-reload",
		"sudo systemctl enable --now " + serviceFile,
	})
}

func TestInstallAndEnable_ReturnsErrorWhenServiceFileMissing(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}
	manifest := &config.ServiceManifest{Entries: []string{"missing.service"}}

	if err := InstallAndEnable(m, t.TempDir(), manifest); err == nil {
		t.Fatalf("expected missing service file error")
	}
}

func TestInstallAndEnable_ReturnsCopyError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()

	dir := t.TempDir()
	serviceFile := "custom.service"
	src := filepath.Join(dir, serviceFile)
	if err := os.WriteFile(src, []byte("[Unit]\nDescription=Custom\n"), 0o644); err != nil {
		t.Fatalf("write service file: %v", err)
	}

	m := &shell.MockExecutor{
		RunErrFor: map[string]error{
			"sudo cp " + src + " /etc/systemd/system/" + serviceFile: errors.New("cp failed"),
		},
	}
	manifest := &config.ServiceManifest{Entries: []string{serviceFile}}

	if err := InstallAndEnable(m, dir, manifest); err == nil {
		t.Fatalf("expected copy error")
	}
}

func assertCalls(t *testing.T, got []shell.CommandCall, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("calls length mismatch\n got: %#v\nwant: %#v", got, want)
	}
	for i := range want {
		gotCall := shell.CommandKey(got[i].Name, got[i].Args...)
		if gotCall != want[i] {
			t.Fatalf("call[%d] mismatch got %q want %q", i, gotCall, want[i])
		}
	}
}
