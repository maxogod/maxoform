package npm

import (
	"fmt"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
)

func BootstrapWithNvm(runner shell.Executor, bootstrap config.NpmBootstrapConfig) error {
	if !bootstrap.Enabled {
		return nil
	}
	if bootstrap.InstallScriptURL == "" {
		return fmt.Errorf("npm bootstrap is enabled but install_script_url is empty")
	}
	if bootstrap.NodeVersion == "" {
		return fmt.Errorf("npm bootstrap is enabled but node_version is empty")
	}

	if err := runner.RunShell(fmt.Sprintf("curl -o- %q | PROFILE=/dev/null bash", bootstrap.InstallScriptURL)); err != nil {
		return err
	}

	return runner.RunShell(fmt.Sprintf("source \"$HOME/.nvm/nvm.sh\" && nvm install %q", bootstrap.NodeVersion))
}

func InstallGlobal(runner shell.Executor, packages []string) error {
	return installGlobal(runner, packages, false)
}

func InstallGlobalWithNvm(runner shell.Executor, packages []string) error {
	return installGlobal(runner, packages, true)
}

func installGlobal(runner shell.Executor, packages []string, useNvm bool) error {
	for _, pkg := range packages {
		if isInstalled(runner, pkg, useNvm) {
			logger.Log.Infof("npm package already installed, skipping: %s", pkg)
			continue
		}

		if useNvm {
			if err := runner.RunShell(fmt.Sprintf("source \"$HOME/.nvm/nvm.sh\" && npm install -g %q", pkg)); err != nil {
				return err
			}
			continue
		}

		if err := runner.Run("npm", "install", "-g", pkg); err != nil {
			return err
		}
	}

	return nil
}

func isInstalled(runner shell.Executor, pkg string, useNvm bool) bool {
	if useNvm {
		return runner.CheckShell(fmt.Sprintf("source \"$HOME/.nvm/nvm.sh\" && npm list -g --depth=0 %q", pkg))
	}

	return runner.Check("npm", "list", "-g", "--depth=0", pkg)
}
