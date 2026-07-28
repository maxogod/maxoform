package app

import (
	"io"
	"testing"

	"github.com/maxogod/maxoform/config"
)

type mockRunner struct {
	shellCalls []string
}

func (m *mockRunner) Run(_ string, _ ...string) error                       { return nil }
func (m *mockRunner) RunWithStdin(_ io.Reader, _ string, _ ...string) error { return nil }
func (m *mockRunner) Check(_ string, _ ...string) bool                      { return false }
func (m *mockRunner) CheckShell(_ string) bool                              { return false }
func (m *mockRunner) RunShell(command string) error {
	m.shellCalls = append(m.shellCalls, command)
	return nil
}

func TestRunPostCommands_RunsCommandsInOrder(t *testing.T) {
	m := &mockRunner{}
	a := &application{
		cfg: &config.Config{
			Commands: config.CommandsConfig{
				Post: []config.ShellCommand{
					{Run: "echo one"},
					{Run: "echo two"},
				},
			},
		},
		runner: m,
	}

	if err := a.runPostCommands(); err != nil {
		t.Fatalf("runPostCommands failed: %v", err)
	}
	if len(m.shellCalls) != 2 || m.shellCalls[0] != "echo one" || m.shellCalls[1] != "echo two" {
		t.Fatalf("unexpected shell calls: %#v", m.shellCalls)
	}
}

func TestRunPostCommands_RejectsEmptyCommand(t *testing.T) {
	a := &application{
		cfg: &config.Config{
			Commands: config.CommandsConfig{
				Post: []config.ShellCommand{{Run: "   "}},
			},
		},
		runner: &mockRunner{},
	}
	if err := a.runPostCommands(); err == nil {
		t.Fatalf("expected empty command error")
	}
}
