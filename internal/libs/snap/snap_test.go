package snap

import (
	"errors"
	"testing"

	"github.com/maxogod/maxoform/internal/logger"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"go.uber.org/zap"
)

func TestRefresh_RunsSnapRefresh(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}
	if err := Refresh(m); err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}
	assertCalls(t, m.RunCalls, []string{"sudo snap refresh"})
}

func TestInstall_SkipsInstalledAndInstallsMissing(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckFor: map[string]bool{
			"snap list nvim": true,
			"snap list vlc":  false,
		},
	}
	if err := Install(m, []string{"nvim", "vlc"}); err != nil {
		t.Fatalf("Install failed: %v", err)
	}
	assertCalls(t, m.RunCalls, []string{"sudo snap install vlc"})
}

func TestInstall_ReturnsInstallError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckFor: map[string]bool{"snap list vlc": false},
		RunErrFor: map[string]error{
			"sudo snap install vlc": errors.New("boom"),
		},
	}
	if err := Install(m, []string{"vlc"}); err == nil {
		t.Fatalf("expected install error")
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
