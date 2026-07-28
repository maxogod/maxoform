package git

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/shell"
)

func TestCloneMissingRepos_SkipsExisting(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "existing-repo")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	m := &shell.MockExecutor{}
	repos := []config.Repo{{URL: "https://github.com/org/repo.git", Dest: dest}}

	if err := CloneMissingRepos(m, repos); err != nil {
		t.Fatalf("CloneMissingRepos failed: %v", err)
	}
	if len(m.RunCalls) != 0 {
		t.Fatalf("expected no run calls for existing repo, got %d", len(m.RunCalls))
	}
}

func TestCloneMissingRepos_ClonesHTTPSAndSetsSSHRemote(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "new-repo")

	m := &shell.MockExecutor{}
	repos := []config.Repo{{URL: "https://github.com/org/repo.git", Dest: dest}}

	if err := CloneMissingRepos(m, repos); err != nil {
		t.Fatalf("CloneMissingRepos failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{
		"git clone https://github.com/org/repo.git " + dest,
		"git -C " + dest + " remote set-url origin git@github.com:org/repo.git",
	})
}

func TestCloneMissingRepos_SkipsSSHRemoteForSSHURL(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "new-repo")

	m := &shell.MockExecutor{}
	repos := []config.Repo{{URL: "git@github.com:org/repo.git", Dest: dest}}

	if err := CloneMissingRepos(m, repos); err != nil {
		t.Fatalf("CloneMissingRepos failed: %v", err)
	}

	// An ssh:// style remote should only trigger the clone, not the
	// remote set-url step.
	assertCalls(t, m.RunCalls, []string{
		"git clone git@github.com:org/repo.git " + dest,
	})
}

func TestCloneMissingRepos_ReturnsCloneError(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "new-repo")

	cloneErr := errors.New("clone failed")
	m := &shell.MockExecutor{
		RunErrFor: map[string]error{
			"git clone https://github.com/org/repo.git " + dest: cloneErr,
		},
	}
	repos := []config.Repo{{URL: "https://github.com/org/repo.git", Dest: dest}}

	if err := CloneMissingRepos(m, repos); err == nil {
		t.Fatalf("expected clone error")
	}
	if len(m.RunCalls) != 1 {
		t.Fatalf("expected clone to stop after error, got %d calls", len(m.RunCalls))
	}
}

func TestCloneMissingRepos_ReturnsSetRemoteError(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "new-repo")

	remoteErr := errors.New("remote set failed")
	m := &shell.MockExecutor{
		RunErrFor: map[string]error{
			"git -C " + dest + " remote set-url origin git@github.com:org/repo.git": remoteErr,
		},
	}
	repos := []config.Repo{{URL: "https://github.com/org/repo.git", Dest: dest}}

	if err := CloneMissingRepos(m, repos); err == nil {
		t.Fatalf("expected remote set-url error")
	}
}

func TestCloneMissingRepos_InvalidRepoURL(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "new-repo")

	m := &shell.MockExecutor{}
	// A host-less https URL makes toSSHRemote return the
	// "invalid repository URL" error.
	repos := []config.Repo{{URL: "https:///no-host.git", Dest: dest}}

	if err := CloneMissingRepos(m, repos); err == nil {
		t.Fatalf("expected error converting remote to ssh")
	}
}

func TestConfigureGlobalUser_SetsNameAndEmail(t *testing.T) {
	m := &shell.MockExecutor{}
	if err := ConfigureGlobalUser(m, "Jane Doe", "jane@example.com"); err != nil {
		t.Fatalf("ConfigureGlobalUser failed: %v", err)
	}
	assertCalls(t, m.RunCalls, []string{
		"git config --global user.name Jane Doe",
		"git config --global user.email jane@example.com",
	})
}

func TestConfigureGlobalUser_SkipsEmptyValues(t *testing.T) {
	m := &shell.MockExecutor{}
	if err := ConfigureGlobalUser(m, "", ""); err != nil {
		t.Fatalf("ConfigureGlobalUser failed: %v", err)
	}
	if len(m.RunCalls) != 0 {
		t.Fatalf("expected no run calls, got %d", len(m.RunCalls))
	}
}

func TestConfigureGlobalUser_NameOnly(t *testing.T) {
	m := &shell.MockExecutor{}
	if err := ConfigureGlobalUser(m, "Jane Doe", ""); err != nil {
		t.Fatalf("ConfigureGlobalUser failed: %v", err)
	}
	assertCalls(t, m.RunCalls, []string{"git config --global user.name Jane Doe"})
}

func TestConfigureGlobalUser_ReturnsRunError(t *testing.T) {
	runErr := errors.New("git failed")
	m := &shell.MockExecutor{
		RunErrFor: map[string]error{
			"git config --global user.name Jane Doe": runErr,
		},
	}
	if err := ConfigureGlobalUser(m, "Jane Doe", "jane@example.com"); err == nil {
		t.Fatalf("expected error from ConfigureGlobalUser")
	}
	if len(m.RunCalls) != 1 {
		t.Fatalf("expected email config to be skipped after name error, got %d calls", len(m.RunCalls))
	}
}

func TestToSSHRemote(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "https url", in: "https://github.com/org/repo.git", want: "git@github.com:org/repo.git"},
		{name: "https url without .git suffix", in: "https://github.com/org/repo", want: "git@github.com:org/repo.git"},
		{name: "http url", in: "http://gitlab.com/org/repo.git", want: "git@gitlab.com:org/repo.git"},
		{name: "already ssh", in: "git@github.com:org/repo.git", want: ""},
		{name: "unsupported scheme returns empty", in: "ftp://example.com/repo.git", want: ""},
		{name: "invalid url missing host", in: "https:///repo.git", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := toSSHRemote(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tc.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("toSSHRemote(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
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
