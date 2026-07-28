package snap

import (
	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
)

func Refresh(runner shell.Executor) error {
	return runner.Run("sudo", "snap", "refresh")
}

func Install(runner shell.Executor, packages []config.SnapPackage) error {
	for _, pkg := range packages {
		args := []string{"snap", "install", pkg.Name}

		if pkg.Classic {
			args = append(args, "--classic")
		}

		if isInstalled(runner, pkg.Name) {
			logger.Log.Infof("snap package already installed, skipping: %s", pkg)
			continue
		}

		if err := runner.Run("sudo", args...); err != nil {
			return err
		}
	}

	return nil
}

func isInstalled(runner shell.Executor, pkg string) bool {
	return runner.Check("snap", "list", pkg)
}
