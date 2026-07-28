package npm

import (
	"fmt"
	"os/exec"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
)

func BootstrapWithNvm(runner shell.Runner, bootstrap config.NpmBootstrapConfig) error {
	if !bootstrap.Enabled {
		return nil
	}
	if bootstrap.InstallScriptURL == "" {
		return fmt.Errorf("npm bootstrap is enabled but install_script_url is empty")
	}
	if bootstrap.NodeVersion == "" {
		return fmt.Errorf("npm bootstrap is enabled but node_version is empty")
	}

	if err := runner.RunShell(fmt.Sprintf("curl -o- %q | bash", bootstrap.InstallScriptURL)); err != nil {
		return err
	}

	return runner.RunShell(fmt.Sprintf("source \"$HOME/.nvm/nvm.sh\" && nvm install %q", bootstrap.NodeVersion))
}

func InstallGlobal(runner shell.Runner, packages []string) error {
	return installGlobal(runner, packages, false)
}

func InstallGlobalWithNvm(runner shell.Runner, packages []string) error {
	return installGlobal(runner, packages, true)
}

func installGlobal(runner shell.Runner, packages []string, useNvm bool) error {
	for _, pkg := range packages {
		if isInstalled(pkg, useNvm) {
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

func isInstalled(pkg string, useNvm bool) bool {
	if useNvm {
		cmd := exec.Command("bash", "-lc", fmt.Sprintf("source \"$HOME/.nvm/nvm.sh\" && npm list -g --depth=0 %q >/dev/null 2>&1", pkg))
		return cmd.Run() == nil
	}

	cmd := exec.Command("npm", "list", "-g", "--depth=0", pkg)
	return cmd.Run() == nil
}
