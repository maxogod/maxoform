package npm

import (
	"os/exec"

	"github.com/maxogod/maxoform/internal/logger"
	"github.com/maxogod/maxoform/internal/libs/shell"
)

func InstallGlobal(runner shell.Runner, packages []string) error {
	for _, pkg := range packages {
		if isInstalled(pkg) {
			logger.Log.Infof("npm package already installed, skipping: %s", pkg)
			continue
		}

		if err := runner.Run("npm", "install", "-g", pkg); err != nil {
			return err
		}
	}

	return nil
}

func isInstalled(pkg string) bool {
	cmd := exec.Command("npm", "list", "-g", "--depth=0", pkg)
	return cmd.Run() == nil
}
