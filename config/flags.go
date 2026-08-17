package config

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/maxogod/maxoform/internal/utils"
)

type FlagConfig struct {
	ConfigPath      string
	SettingsDirPath string
	ServicesDirPath string

	DconfManifestPath    string
	ServicesManifestPath string

	Quiet    bool
	QuietAll bool

	InstallsOnly bool
	ReposOnly    bool
	SettingsOnly bool
	CommandsOnly bool
	ServicesOnly bool
}

func LoadFlagConf() (*FlagConfig, error) {
	cfgPath := flag.String("config", "", "path to packages YAML")
	settingsDir := flag.String("settings-dir", "", "path to dconf .ini files directory")
	servicesDir := flag.String("services-dir", "", "path to service files directory")

	quiet := flag.Bool("quiet", false, "suppress output from executed commands, keep application logs")
	quietAll := flag.Bool("quiet-all", false, "suppress all output, including application logs")

	installsOnly := flag.Bool("installs-only", false, "only run package installs (apt, snap, npm, pipx)")
	reposOnly := flag.Bool("repos-only", false, "only clone/update repos and set up ssh")
	settingsOnly := flag.Bool("settings-only", false, "only apply dconf settings")
	commandsOnly := flag.Bool("commands-only", false, "only run post-install commands")
	servicesOnly := flag.Bool("services-only", false, "only set up systemd services")
	flag.Parse()

	flagCfg := &FlagConfig{
		ConfigPath:      *cfgPath,
		SettingsDirPath: *settingsDir,
		ServicesDirPath: *servicesDir,

		Quiet:    *quiet,
		QuietAll: *quietAll,

		InstallsOnly: *installsOnly,
		ReposOnly:    *reposOnly,
		SettingsOnly: *settingsOnly,
		CommandsOnly: *commandsOnly,
		ServicesOnly: *servicesOnly,
	}

	if err := validateFlags(flagCfg); err != nil {
		flag.Usage()
		return nil, err
	}

	return flagCfg, nil
}

// anyStepOnlySet reports whether at least one "-only" flag was passed. When
// none are set, every step runs (the default, unrestricted behavior).
func (f *FlagConfig) anyStepOnlySet() bool {
	return f.InstallsOnly || f.ReposOnly || f.SettingsOnly || f.CommandsOnly || f.ServicesOnly
}

// RunInstalls reports whether the apt/snap/npm/pipx install step should run.
func (f *FlagConfig) RunInstalls() bool {
	return !f.anyStepOnlySet() || f.InstallsOnly
}

// RunRepos reports whether the repo clone/update + ssh key step should run.
func (f *FlagConfig) RunRepos() bool {
	return !f.anyStepOnlySet() || f.ReposOnly
}

// RunSettings reports whether the dconf settings import step should run.
func (f *FlagConfig) RunSettings() bool {
	return !f.anyStepOnlySet() || f.SettingsOnly
}

// RunCommands reports whether the post-install commands step should run.
func (f *FlagConfig) RunCommands() bool {
	return !f.anyStepOnlySet() || f.CommandsOnly
}

// RunServices reports whether the systemd services step should run.
func (f *FlagConfig) RunServices() bool {
	return !f.anyStepOnlySet() || f.ServicesOnly
}

func validateFlags(flagCfg *FlagConfig) error {
	if flagCfg.ConfigPath == "" || !utils.PathExists(flagCfg.ConfigPath) {
		return fmt.Errorf("missing --config")
	}
	if flagCfg.SettingsDirPath == "" || !utils.PathExists(flagCfg.SettingsDirPath) {
		return fmt.Errorf("missing --settings-dir")
	}
	if flagCfg.ServicesDirPath == "" || !utils.PathExists(flagCfg.ServicesDirPath) {
		return fmt.Errorf("missing --services-dir")
	}

	flagCfg.DconfManifestPath = filepath.Join(flagCfg.SettingsDirPath, "manifest.yaml")
	if !utils.PathExists(flagCfg.DconfManifestPath) {
		return fmt.Errorf("missing settings manifest at %s", flagCfg.DconfManifestPath)
	}

	flagCfg.ServicesManifestPath = filepath.Join(flagCfg.ServicesDirPath, "manifest.yaml")
	if !utils.PathExists(flagCfg.ServicesManifestPath) {
		return fmt.Errorf("missing services manifest at %s", flagCfg.ServicesManifestPath)
	}
	return nil
}
