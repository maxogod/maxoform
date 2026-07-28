package npm

import (
	"errors"
	"testing"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/logger"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"go.uber.org/zap"
)

func TestBootstrapWithNvm_Validation(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}
	if err := BootstrapWithNvm(m, config.NpmBootstrapConfig{Enabled: true, NodeVersion: "26"}); err == nil {
		t.Fatalf("expected missing script url error")
	}
	if err := BootstrapWithNvm(m, config.NpmBootstrapConfig{Enabled: true, InstallScriptURL: "https://x"}); err == nil {
		t.Fatalf("expected missing node version error")
	}
}

func TestBootstrapWithNvm_RunsCommands(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{}
	err := BootstrapWithNvm(m, config.NpmBootstrapConfig{
		Enabled:          true,
		InstallScriptURL: "https://example.com/install.sh",
		NodeVersion:      "26",
	})
	if err != nil {
		t.Fatalf("BootstrapWithNvm failed: %v", err)
	}
	if len(m.RunShellCalls) != 2 {
		t.Fatalf("expected 2 shell calls, got %d", len(m.RunShellCalls))
	}
}

func TestInstallGlobal_UsesBinaryChecksAndRun(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckFor: map[string]bool{
			"npm list -g --depth=0 eslint":     true,
			"npm list -g --depth=0 typescript": false,
		},
	}
	if err := InstallGlobal(m, []string{"eslint", "typescript"}); err != nil {
		t.Fatalf("InstallGlobal failed: %v", err)
	}
	assertCalls(t, m.RunCalls, []string{"npm install -g typescript"})
}

func TestInstallGlobalWithNvm_UsesShellChecksAndShellInstall(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	checkCmdA := "source \"$HOME/.nvm/nvm.sh\" && npm list -g --depth=0 \"eslint\""
	checkCmdB := "source \"$HOME/.nvm/nvm.sh\" && npm list -g --depth=0 \"typescript\""
	installCmdB := "source \"$HOME/.nvm/nvm.sh\" && npm install -g \"typescript\""
	m := &shell.MockExecutor{
		CheckShellFor: map[string]bool{
			checkCmdA: true,
			checkCmdB: false,
		},
	}
	if err := InstallGlobalWithNvm(m, []string{"eslint", "typescript"}); err != nil {
		t.Fatalf("InstallGlobalWithNvm failed: %v", err)
	}
	assertShellCalls(t, m.RunShellCalls, []string{installCmdB})
}

func TestInstallGlobal_ReturnsInstallError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()
	m := &shell.MockExecutor{
		CheckFor: map[string]bool{"npm list -g --depth=0 typescript": false},
		RunErrFor: map[string]error{
			"npm install -g typescript": errors.New("boom"),
		},
	}
	if err := InstallGlobal(m, []string{"typescript"}); err == nil {
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

func assertShellCalls(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("calls length mismatch\n got: %#v\nwant: %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("call[%d] mismatch got %q want %q", i, got[i], want[i])
		}
	}
}
