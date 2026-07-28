package apt

import (
	"errors"
	"testing"

	"github.com/maxogod/maxoform/internal/logger"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"go.uber.org/zap"
)

func TestUpdateSystem_RunsExpectedCommandsInOrder(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}

	if err := UpdateSystem(m); err != nil {
		t.Fatalf("UpdateSystem failed: %v", err)
	}

	want := []string{
		"sudo apt update",
		"sudo apt upgrade -y",
		"sudo apt autoremove -y",
	}
	assertCalls(t, m.RunCalls, want)
}

func TestInstall_SkipsInstalledAndInstallsMissing(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckFor: map[string]bool{
			"dpkg -s git":  true,
			"dpkg -s curl": false,
		},
	}

	if err := Install(m, []string{"git", "curl"}); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{"sudo apt install -y curl"})
}

func TestInstall_ReturnsInstallError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckFor: map[string]bool{"dpkg -s curl": false},
		RunErrFor: map[string]error{
			"sudo apt install -y curl": errors.New("boom"),
		},
	}

	if err := Install(m, []string{"curl"}); err == nil {
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
