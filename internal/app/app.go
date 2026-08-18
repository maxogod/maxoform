package app

import (
	"fmt"
	"strings"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/apt"
	"github.com/maxogod/maxoform/internal/libs/dconf"
	"github.com/maxogod/maxoform/internal/libs/git"
	"github.com/maxogod/maxoform/internal/libs/npm"
	"github.com/maxogod/maxoform/internal/libs/pipx"
	"github.com/maxogod/maxoform/internal/libs/services"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/libs/snap"
	"github.com/maxogod/maxoform/internal/libs/sshkey"
	"github.com/maxogod/maxoform/internal/libs/swap"
	"github.com/maxogod/maxoform/internal/logger"
)

type application struct {
	cfg             *config.Config
	flagCfg         *config.FlagConfig
	dconfManifest   *config.DconfManifest
	serviceManifest *config.ServiceManifest
	runner          shell.Executor
}

func New(cfg *config.Config, flagCfg *config.FlagConfig, dconfManifest *config.DconfManifest, serviceManifest *config.ServiceManifest) Application {
	return &application{
		cfg:             cfg,
		flagCfg:         flagCfg,
		dconfManifest:   dconfManifest,
		serviceManifest: serviceManifest,
		runner:          shell.Runner{Quiet: flagCfg.Quiet || flagCfg.QuietAll},
	}
}

func (a *application) Run() error {
	if a.flagCfg.RunInstalls() {
		logger.Log.Info("1/7: Updating package managers")
		if err := apt.UpdateSystem(a.runner); err != nil {
			return fmt.Errorf("updating apt packages: %w", err)
		}
		if err := snap.Refresh(a.runner); err != nil {
			return fmt.Errorf("refreshing snap packages: %w", err)
		}

		logger.Log.Info("2/7: Installing packages")
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
		if err := pipx.Install(a.runner, a.cfg.Packages.Pipx); err != nil {
			return fmt.Errorf("installing pipx packages: %w", err)
		}
	} else {
		logger.Log.Info("1-2/7: Skipping package manager updates and installs")
	}

	if a.flagCfg.RunRepos() {
		logger.Log.Info("3/7: Cloning repositories")
		if err := git.ConfigureGlobalUser(a.runner, a.cfg.Settings.GitUserName, a.cfg.Settings.GitUserEmail, a.cfg.Settings.GitCoreEditor); err != nil {
			return fmt.Errorf("configuring git user: %w", err)
		}
		if err := git.CloneMissingRepos(a.runner, a.cfg.Repos); err != nil {
			return fmt.Errorf("cloning repositories: %w", err)
		}
	} else {
		logger.Log.Info("3/7: Skipping repository cloning")
	}

	if a.flagCfg.RunSettings() {
		if a.cfg.Settings.ImportDconf {
			logger.Log.Info("4/7: Importing dconf settings")
			if err := dconf.Import(a.runner, a.flagCfg.SettingsDirPath, a.dconfManifest); err != nil {
				return fmt.Errorf("importing dconf settings: %w", err)
			}
		} else {
			logger.Log.Info("4/7: Skipping GNOME dconf settings import")
		}

		logger.Log.Info("4/7: Configuring swap file")
		if err := swap.Ensure(a.runner, a.cfg.Settings.SwapSizeGB); err != nil {
			return fmt.Errorf("configuring swap file: %w", err)
		}
	} else {
		logger.Log.Info("4/7: Skipping GNOME dconf settings import")
		logger.Log.Info("4/7: Skipping swap file configuration")
	}

	if a.flagCfg.RunRepos() {
		logger.Log.Info("5/7: Creating SSH key")
		if err := sshkey.EnsureAndPrint(a.runner, a.cfg.Settings.GitUserEmail, a.cfg.Settings.SshKeyPath, a.cfg.Settings.SshPassphrase); err != nil {
			return fmt.Errorf("creating ssh key: %w", err)
		}
	} else {
		logger.Log.Info("5/7: Skipping SSH key creation")
	}

	if a.flagCfg.RunCommands() {
		logger.Log.Info("6/7: Running post commands")
		if err := a.runPostCommands(); err != nil {
			return err
		}
	} else {
		logger.Log.Info("6/7: Skipping post commands")
	}

	if a.flagCfg.RunServices() {
		logger.Log.Info("7/7: Installing and enabling systemd services")
		if err := services.InstallAndEnable(a.runner, a.flagCfg.ServicesDirPath, a.serviceManifest); err != nil {
			return fmt.Errorf("installing and enabling services: %w", err)
		}
	} else {
		logger.Log.Info("7/7: Skipping systemd services setup")
	}

	logger.Log.Info("maxoform completed. You should now reboot the system for some settings to apply.")
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
