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
		Entries: []config.ServiceEntry{
			{Condition: "test -f /tmp/exists", Service: serviceFile},
		},
	}
	m.CheckShellFor = map[string]bool{"test -f /tmp/exists": true}

	if err := InstallAndEnable(m, dir, manifest); err != nil {
		t.Fatalf("InstallAndEnable failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{
		"sudo cp " + filepath.Join(dir, serviceFile) + " /etc/systemd/system/" + serviceFile,
		"sudo systemctl daemon-reload",
		"sudo systemctl enable --now " + serviceFile,
	})
}

func TestInstallAndEnable_SkipsWhenConditionFails(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}
	manifest := &config.ServiceManifest{
		Entries: []config.ServiceEntry{
			{Condition: "test -f /tmp/missing", Service: "custom.service"},
		},
	}
	m.CheckShellFor = map[string]bool{"test -f /tmp/missing": false}

	if err := InstallAndEnable(m, t.TempDir(), manifest); err != nil {
		t.Fatalf("expected no error when condition is false, got: %v", err)
	}
	if len(m.RunCalls) != 0 {
		t.Fatalf("expected no run calls when condition is false, got %#v", m.RunCalls)
	}
}

func TestInstallAndEnable_ReturnsErrorWhenServiceFileMissing(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}
	manifest := &config.ServiceManifest{
		Entries: []config.ServiceEntry{
			{Condition: "true", Service: "missing.service"},
		},
	}
	m.CheckShellFor = map[string]bool{"true": true}

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
	manifest := &config.ServiceManifest{
		Entries: []config.ServiceEntry{
			{Condition: "true", Service: serviceFile},
		},
	}
	m.CheckShellFor = map[string]bool{"true": true}

	if err := InstallAndEnable(m, dir, manifest); err == nil {
		t.Fatalf("expected copy error")
	}
}

func TestInstallAndEnable_ReturnsErrorWhenServiceFieldEmpty(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}
	manifest := &config.ServiceManifest{
		Entries: []config.ServiceEntry{
			{Condition: "true", Service: "   "},
		},
	}
	if err := InstallAndEnable(m, t.TempDir(), manifest); err == nil {
		t.Fatalf("expected error for empty service")
	}
}

func TestInstallAndEnable_ReturnsErrorWhenConditionFieldEmpty(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}
	manifest := &config.ServiceManifest{
		Entries: []config.ServiceEntry{
			{Condition: "   ", Service: "custom.service"},
		},
	}
	if err := InstallAndEnable(m, t.TempDir(), manifest); err == nil {
		t.Fatalf("expected error for empty condition")
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
