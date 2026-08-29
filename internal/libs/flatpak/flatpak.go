package flatpak

import (
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
)

const (
	remoteName = "flathub"
	remoteURL  = "https://dl.flathub.org/repo/flathub.flatpakrepo"
)

// AddFlathubRemote registers the flathub remote system-wide. It is a no-op when
// the remote already exists.
func AddFlathubRemote(runner shell.Executor) error {
	return runner.Run("flatpak", "remote-add", "--if-not-exists", remoteName, remoteURL)
}

// Install installs the given flatpak application ids from flathub, skipping the
// ones already present.
func Install(runner shell.Executor, packages []string) error {
	for _, pkg := range packages {
		if isInstalled(runner, pkg) {
			logger.Log.Infof("flatpak package already installed, skipping: %s", pkg)
			continue
		}

		if err := runner.Run("flatpak", "install", remoteName, pkg, "-y"); err != nil {
			return err
		}
	}

	return nil
}

func isInstalled(runner shell.Executor, pkg string) bool {
	return runner.Check("flatpak", "info", pkg)
}
