package snap

import (
	"os/exec"

	"github.com/maxogod/maxoform/internal/logger"
	"github.com/maxogod/maxoform/internal/libs/shell"
)

func Refresh(runner shell.Runner) error {
	return runner.Run("sudo", "snap", "refresh")
}

func Install(runner shell.Runner, packages []string) error {
	for _, pkg := range packages {
		if isInstalled(pkg) {
			logger.Log.Infof("snap package already installed, skipping: %s", pkg)
			continue
		}

		if err := runner.Run("sudo", "snap", "install", pkg); err != nil {
			return err
		}
	}

	return nil
}

func isInstalled(pkg string) bool {
	cmd := exec.Command("snap", "list", pkg)
	return cmd.Run() == nil
}
