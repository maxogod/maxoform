package sshkey

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/maxogod/maxoform/internal/libs/shell"
)

func TestEnsureAndPrint_RequiresEmail(t *testing.T) {
	m := &shell.MockExecutor{}
	if err := EnsureAndPrint(m, "", "/tmp/id_ed25519", ""); err == nil {
		t.Fatalf("expected error for missing email")
	}
	if len(m.RunCalls) != 0 {
		t.Fatalf("expected no run calls when email missing, got %d", len(m.RunCalls))
	}
}

func TestEnsureAndPrint_GeneratesKeyWhenMissing(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "id_ed25519")
	pubPath := keyPath + ".pub"

	m := &shell.MockExecutor{
		// ssh-keygen isn't actually run in tests, so simulate the file
		// it would have produced.
		OnRun: func(name string, args ...string) error {
			return os.WriteFile(pubPath, []byte("ssh-ed25519 AAAA... test@example.com\n"), 0o644)
		},
	}

	if err := EnsureAndPrint(m, "test@example.com", keyPath, "passphrase"); err != nil {
		t.Fatalf("EnsureAndPrint failed: %v", err)
	}

	assertCalls(t, m.RunCalls, []string{
		"ssh-keygen -t ed25519 -C test@example.com -f " + keyPath + " -N passphrase",
	})
}

func TestEnsureAndPrint_SkipsGenerationWhenKeyExists(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "id_ed25519")
	pubPath := keyPath + ".pub"

	if err := os.WriteFile(pubPath, []byte("ssh-ed25519 AAAA... existing@example.com\n"), 0o644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	m := &shell.MockExecutor{}
	if err := EnsureAndPrint(m, "existing@example.com", keyPath, ""); err != nil {
		t.Fatalf("EnsureAndPrint failed: %v", err)
	}
	if len(m.RunCalls) != 0 {
		t.Fatalf("expected no ssh-keygen call when key already exists, got %d", len(m.RunCalls))
	}
}

func TestEnsureAndPrint_ReturnsKeygenError(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "id_ed25519")

	keygenErr := errors.New("keygen failed")
	m := &shell.MockExecutor{
		RunErrFor: map[string]error{
			"ssh-keygen -t ed25519 -C test@example.com -f " + keyPath + " -N ": keygenErr,
		},
	}

	if err := EnsureAndPrint(m, "test@example.com", keyPath, ""); err == nil {
		t.Fatalf("expected keygen error")
	}
}

func TestEnsureAndPrint_ReturnsReadFileErrorWhenPubKeyMissingAfterGen(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "id_ed25519")

	// OnRun succeeds but never writes the .pub file, simulating a
	// ssh-keygen invocation that didn't produce the expected output.
	m := &shell.MockExecutor{}

	if err := EnsureAndPrint(m, "test@example.com", keyPath, ""); err == nil {
		t.Fatalf("expected error reading missing public key file")
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
