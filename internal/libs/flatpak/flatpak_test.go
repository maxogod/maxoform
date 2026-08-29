package flatpak

import (
	"errors"
	"testing"

	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
	"go.uber.org/zap"
)

func TestAddFlathubRemote_RunsRemoteAdd(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()

	m := &shell.MockExecutor{}

	if err := AddFlathubRemote(m); err != nil {
		t.Fatalf("AddFlathubRemote failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{
		"flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo",
	})
}

func TestAddFlathubRemote_ReturnsRemoteAddError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()

	m := &shell.MockExecutor{
		RunErrFor: map[string]error{
			"flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo": errors.New("boom"),
		},
	}

	if err := AddFlathubRemote(m); err == nil {
		t.Fatalf("expected remote-add error")
	}
}

func TestInstall_SkipsInstalledAndInstallsMissing(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()

	m := &shell.MockExecutor{
		CheckFor: map[string]bool{
			"flatpak info org.kde.krita":      true,
			"flatpak info com.spotify.Client": false,
		},
	}

	if err := Install(m, []string{"org.kde.krita", "com.spotify.Client"}); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{
		"flatpak install flathub com.spotify.Client -y",
	})
}

func TestInstall_ReturnsInstallError(t *testing.T) {
	logger.Log = zap.NewNop().Sugar()

	m := &shell.MockExecutor{
		CheckFor: map[string]bool{
			"flatpak info org.kde.krita": false,
		},
		RunErrFor: map[string]error{
			"flatpak install flathub org.kde.krita -y": errors.New("boom"),
		},
	}

	if err := Install(m, []string{"org.kde.krita"}); err == nil {
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
