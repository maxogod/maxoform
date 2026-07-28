package shell

import (
	"io"
	"strings"
)

type CommandCall struct {
	Name string
	Args []string
}

type MockExecutor struct {
	RunCalls          []CommandCall
	RunWithStdinCalls []CommandCall
	RunShellCalls     []string
	CheckCalls        []CommandCall
	CheckShellCalls   []string
	LastStdinData     string

	RunErrFor      map[string]error
	RunShellErrFor map[string]error
	CheckFor       map[string]bool
	CheckShellFor  map[string]bool

	OnRun          func(name string, args ...string) error
	OnRunWithStdin func(stdin []byte, name string, args ...string) error
}

func (m *MockExecutor) Run(name string, args ...string) error {
	m.RunCalls = append(m.RunCalls, CommandCall{Name: name, Args: append([]string{}, args...)})
	if m.OnRun != nil {
		if err := m.OnRun(name, args...); err != nil {
			return err
		}
	}
	if err, ok := m.RunErrFor[CommandKey(name, args...)]; ok {
		return err
	}
	return nil
}

func (m *MockExecutor) RunWithStdin(stdin io.Reader, name string, args ...string) error {
	data, _ := io.ReadAll(stdin)
	m.LastStdinData = string(data)
	m.RunWithStdinCalls = append(m.RunWithStdinCalls, CommandCall{Name: name, Args: append([]string{}, args...)})
	if m.OnRunWithStdin != nil {
		if err := m.OnRunWithStdin(data, name, args...); err != nil {
			return err
		}
	}
	if err, ok := m.RunErrFor[CommandKey(name, args...)]; ok {
		return err
	}
	return nil
}

func (m *MockExecutor) RunShell(command string) error {
	m.RunShellCalls = append(m.RunShellCalls, command)
	if err, ok := m.RunShellErrFor[command]; ok {
		return err
	}
	return nil
}

func (m *MockExecutor) Check(name string, args ...string) bool {
	m.CheckCalls = append(m.CheckCalls, CommandCall{Name: name, Args: append([]string{}, args...)})
	return m.CheckFor[CommandKey(name, args...)]
}

func (m *MockExecutor) CheckShell(command string) bool {
	m.CheckShellCalls = append(m.CheckShellCalls, command)
	return m.CheckShellFor[command]
}

func CommandKey(name string, args ...string) string {
	if len(args) == 0 {
		return name
	}
	return name + " " + strings.Join(args, " ")
}
