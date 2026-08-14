package pipx

import (
	"errors"
	"testing"

	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
	"go.uber.org/zap"
)

func TestInstall_SkipsInstalledAndInstallsMissing(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()

	checkCmdA := "pipx list --short | grep -Fx \"pre-commit\""
	checkCmdB := "pipx list --short | grep -Fx \"poetry\""
	m := &shell.MockExecutor{
		CheckShellFor: map[string]bool{
			checkCmdA: true,
			checkCmdB: false,
		},
	}

	if err := Install(m, []string{"pre-commit", "poetry"}); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{"pipx install poetry"})
}

func TestInstall_ReturnsInstallError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()

	m := &shell.MockExecutor{
		CheckShellFor: map[string]bool{
			"pipx list --short | grep -Fx \"poetry\"": false,
		},
		RunErrFor: map[string]error{
			"pipx install poetry": errors.New("boom"),
		},
	}

	if err := Install(m, []string{"poetry"}); err == nil {
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
