package shell

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/maxogod/maxoform/internal/logger"
)

type Runner struct {
	// Quiet discards the executed commands' own stdout/stderr while still
	// allowing application logs (e.g. the "$ command" line) to print.
	Quiet bool
}

func (r Runner) stdoutStderr() (io.Writer, io.Writer) {
	if r.Quiet {
		return io.Discard, io.Discard
	}
	return os.Stdout, os.Stderr
}

func (r Runner) Run(name string, args ...string) error {
	all := append([]string{name}, args...)
	logger.Log.Infof("$ %s", strings.Join(all, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = r.stdoutStderr()
	return cmd.Run()
}

func (r Runner) RunWithStdin(stdin io.Reader, name string, args ...string) error {
	all := append([]string{name}, args...)
	logger.Log.Infof("$ %s", strings.Join(all, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout, cmd.Stderr = r.stdoutStderr()
	return cmd.Run()
}

func (r Runner) RunShell(command string) error {
	logger.Log.Infof("$ bash -lc %q", command)

	cmd := exec.Command("bash", "-lc", command)
	cmd.Stdout, cmd.Stderr = r.stdoutStderr()
	return cmd.Run()
}

func (r Runner) Check(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run() == nil
}

func (r Runner) CheckShell(command string) bool {
	cmd := exec.Command("bash", "-lc", command)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run() == nil
}
