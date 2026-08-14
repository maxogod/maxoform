package pipx

import (
	"fmt"

	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
)

func Install(runner shell.Executor, packages []string) error {
	for _, pkg := range packages {
		if isInstalled(runner, pkg) {
			logger.Log.Infof("pipx package already installed, skipping: %s", pkg)
			continue
		}

		if err := runner.Run("pipx", "install", pkg); err != nil {
			return err
		}
	}

	return nil
}

func isInstalled(runner shell.Executor, pkg string) bool {
	return runner.CheckShell(fmt.Sprintf("pipx list --short | grep -Fx %q", pkg))
}
