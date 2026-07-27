package shell

import (
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
