package app

import (
	"fmt"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/apt"
	"github.com/maxogod/maxoform/internal/libs/dconf"
	"github.com/maxogod/maxoform/internal/libs/git"
	"github.com/maxogod/maxoform/internal/libs/npm"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/libs/snap"
	"github.com/maxogod/maxoform/internal/libs/sshkey"
	"github.com/maxogod/maxoform/internal/logger"
)

type application struct {
	cfg      *config.Config
	flagCfg  *config.FlagConfig
	manifest *config.DconfManifest
	runner   shell.Runner
}

func New(cfg *config.Config, flagCfg *config.FlagConfig, manifest *config.DconfManifest) Application {
	return &application{
		cfg:      cfg,
		flagCfg:  flagCfg,
		manifest: manifest,
		runner:   shell.Runner{},
	}
}

func (a *application) Run() error {
	logger.Log.Info("1/5: Updating package managers")
	if err := apt.UpdateSystem(a.runner); err != nil {
		return fmt.Errorf("updating apt packages: %w", err)
	}
	if err := snap.Refresh(a.runner); err != nil {
		return fmt.Errorf("refreshing snap packages: %w", err)
	}

	logger.Log.Info("2/5: Installing packages")
	if err := apt.Install(a.runner, a.cfg.Packages.Apt); err != nil {
		return fmt.Errorf("installing apt packages: %w", err)
	}
	if err := snap.Install(a.runner, a.cfg.Packages.Snap); err != nil {
		return fmt.Errorf("installing snap packages: %w", err)
	}
	if err := npm.InstallGlobal(a.runner, a.cfg.Packages.Npm); err != nil {
		return fmt.Errorf("installing npm packages: %w", err)
	}

	logger.Log.Info("3/5: Cloning repositories")
	if err := git.ConfigureGlobalUser(a.runner, a.cfg.Settings.GitUserName, a.cfg.Settings.GitUserEmail); err != nil {
		return fmt.Errorf("configuring git user: %w", err)
	}
	if err := git.CloneMissingRepos(a.runner, a.cfg.Repos); err != nil {
		return fmt.Errorf("cloning repositories: %w", err)
	}

	if a.cfg.Settings.ImportDconf {
		logger.Log.Info("4/5: Importing dconf settings")
		if err := dconf.Import(a.runner, a.flagCfg.SettingsDirPath, a.manifest); err != nil {
			return fmt.Errorf("importing dconf settings: %w", err)
		}
	} else {
		logger.Log.Info("4/5: Skipping GNOME dconf settings import")
	}

	logger.Log.Info("5/5: Creating SSH key")
	if err := sshkey.EnsureAndPrint(a.runner, a.cfg.Settings.GitUserEmail, a.cfg.Settings.SshKeyPath, a.cfg.Settings.SshPassphrase); err != nil {
		return fmt.Errorf("creating ssh key: %w", err)
	}

	logger.Log.Info("maxoform completed")
}
