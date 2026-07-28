package shell

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/maxogod/maxoform/internal/logger"
)

type Runner struct {
}

func (r Runner) Run(name string, args ...string) error {
	all := append([]string{name}, args...)
	logger.Log.Infof("$ %s", strings.Join(all, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r Runner) RunWithStdin(stdin io.Reader, name string, args ...string) error {
	all := append([]string{name}, args...)
	logger.Log.Infof("$ %s", strings.Join(all, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r Runner) RunShell(command string) error {
	logger.Log.Infof("$ bash -lc %q", command)

	cmd := exec.Command("bash", "-lc", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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
