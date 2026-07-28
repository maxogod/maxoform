package shell

import "io"

type Executor interface {
	Run(name string, args ...string) error
	RunWithStdin(stdin io.Reader, name string, args ...string) error
	RunShell(command string) error
	Check(name string, args ...string) bool
	CheckShell(command string) bool
}
