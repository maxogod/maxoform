package app

import (
	"fmt"
	"strings"

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
	logger.Log.Info("1/6: Updating package managers")
	if err := apt.UpdateSystem(a.runner); err != nil {
		return fmt.Errorf("updating apt packages: %w", err)
	}
	if err := snap.Refresh(a.runner); err != nil {
		return fmt.Errorf("refreshing snap packages: %w", err)
	}

	logger.Log.Info("2/6: Installing packages")
	if err := apt.Install(a.runner, a.cfg.Packages.Apt); err != nil {
		return fmt.Errorf("installing apt packages: %w", err)
	}
	if err := snap.Install(a.runner, a.cfg.Packages.Snap); err != nil {
		return fmt.Errorf("installing snap packages: %w", err)
	}

	if a.cfg.NpmBootstrap.Enabled {
		if err := npm.BootstrapWithNvm(a.runner, a.cfg.NpmBootstrap); err != nil {
			return fmt.Errorf("bootstrapping npm with nvm: %w", err)
		}
		if err := npm.InstallGlobalWithNvm(a.runner, a.cfg.Packages.Npm); err != nil {
			return fmt.Errorf("installing npm packages: %w", err)
		}
	} else if err := npm.InstallGlobal(a.runner, a.cfg.Packages.Npm); err != nil {
		return fmt.Errorf("installing npm packages: %w", err)
	}

	logger.Log.Info("3/6: Cloning repositories")
	if err := git.ConfigureGlobalUser(a.runner, a.cfg.Settings.GitUserName, a.cfg.Settings.GitUserEmail); err != nil {
		return fmt.Errorf("configuring git user: %w", err)
	}
	if err := git.CloneMissingRepos(a.runner, a.cfg.Repos); err != nil {
		return fmt.Errorf("cloning repositories: %w", err)
	}

	if a.cfg.Settings.ImportDconf {
		logger.Log.Info("4/6: Importing dconf settings")
		if err := dconf.Import(a.runner, a.flagCfg.SettingsDirPath, a.manifest); err != nil {
			return fmt.Errorf("importing dconf settings: %w", err)
		}
	} else {
		logger.Log.Info("4/6: Skipping GNOME dconf settings import")
	}

	logger.Log.Info("5/6: Creating SSH key")
	if err := sshkey.EnsureAndPrint(a.runner, a.cfg.Settings.GitUserEmail, a.cfg.Settings.SshKeyPath, a.cfg.Settings.SshPassphrase); err != nil {
		return fmt.Errorf("creating ssh key: %w", err)
	}

	logger.Log.Info("6/6: Running post commands")
	if err := a.runPostCommands(); err != nil {
		return err
	}

	logger.Log.Info("maxoform completed")
	return nil
}

func (a *application) runPostCommands() error {
	for idx, cmd := range a.cfg.Commands.Post {
		if strings.TrimSpace(cmd.Run) == "" {
			return fmt.Errorf("running post commands: post command at index %d is empty", idx)
		}
		if err := a.runner.RunShell(cmd.Run); err != nil {
			return fmt.Errorf("running post commands: command at index %d failed: %w", idx, err)
		}
	}

	return nil
}
