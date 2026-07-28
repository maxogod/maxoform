package apt

import (
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
)

func UpdateSystem(runner shell.Runner) error {
	if err := runner.Run("sudo", "apt", "update"); err != nil {
		return err
	}

	if err := runner.Run("sudo", "apt", "upgrade", "-y"); err != nil {
		return err
	}

	return runner.Run("sudo", "apt", "autoremove", "-y")
}

func Install(runner shell.Runner, packages []string) error {
	for _, pkg := range packages {
		if isInstalled(runner, pkg) {
			logger.Log.Infof("apt package already installed, skipping: %s", pkg)
			continue
		}

		if err := runner.Run("sudo", "apt", "install", "-y", pkg); err != nil {
			return err
		}
	}

	return nil
}

func isInstalled(runner shell.Runner, pkg string) bool {
	return runner.Check("dpkg", "-s", pkg)
}
