package swap

import (
	"errors"
	"testing"

	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
	"go.uber.org/zap"
)

func TestEnsure_SkipsWhenSizeNotConfigured(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}

	if err := Ensure(m, 0); err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}
	if len(m.RunCalls) != 0 || len(m.RunShellCalls) != 0 || len(m.CheckShellCalls) != 0 {
		t.Fatalf("expected no calls when size is not configured, got run=%#v runShell=%#v checkShell=%#v", m.RunCalls, m.RunShellCalls, m.CheckShellCalls)
	}
}

func TestEnsure_SkipsWhenAlreadyActiveWithRequestedSize(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckShellFor: map[string]bool{
			"swapon --show=NAME --noheadings | grep -qx /swapfile":    true,
			`[ "$(stat -c%s /swapfile 2>/dev/null)" = "8589934592" ]`: true,
		},
	}

	if err := Ensure(m, 8); err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}
	if len(m.RunCalls) != 0 {
		t.Fatalf("expected no run calls when swap file already matches, got %#v", m.RunCalls)
	}
}

func TestEnsure_CreatesSwapFileWhenNotActive(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckShellFor: map[string]bool{
			"swapon --show=NAME --noheadings | grep -qx /swapfile": false,
			"grep -qF /swapfile /etc/fstab":                        false,
		},
	}

	if err := Ensure(m, 8); err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{
		"sudo fallocate -l 8G /swapfile",
		"sudo chmod 600 /swapfile",
		"sudo mkswap /swapfile",
		"sudo swapon /swapfile",
	})

	if len(m.RunShellCalls) != 1 || m.RunShellCalls[0] != "echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab >/dev/null" {
		t.Fatalf("expected fstab entry to be appended, got %#v", m.RunShellCalls)
	}
}

func TestEnsure_RecreatesWhenActiveWithWrongSize(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckShellFor: map[string]bool{
			"swapon --show=NAME --noheadings | grep -qx /swapfile":    true,
			`[ "$(stat -c%s /swapfile 2>/dev/null)" = "8589934592" ]`: false,
			"grep -qF /swapfile /etc/fstab":                           true,
		},
	}

	if err := Ensure(m, 8); err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{
		"sudo swapoff /swapfile",
		"sudo fallocate -l 8G /swapfile",
		"sudo chmod 600 /swapfile",
		"sudo mkswap /swapfile",
		"sudo swapon /swapfile",
	})

	if len(m.RunShellCalls) != 0 {
		t.Fatalf("expected no fstab entry when already present, got %#v", m.RunShellCalls)
	}
}

func TestEnsure_ReturnsFallocateError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckShellFor: map[string]bool{
			"swapon --show=NAME --noheadings | grep -qx /swapfile": false,
		},
		RunErrFor: map[string]error{
			"sudo fallocate -l 8G /swapfile": errors.New("boom"),
		},
	}

	if err := Ensure(m, 8); err == nil {
		t.Fatalf("expected fallocate error")
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
