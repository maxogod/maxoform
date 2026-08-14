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
}

func LoadFlagConf() (*FlagConfig, error) {
	cfgPath := flag.String("config", "", "path to packages YAML")
	settingsDir := flag.String("settings-dir", "", "path to dconf .ini files directory")
	servicesDir := flag.String("services-dir", "", "path to service files directory")
	flag.Parse()

	flagCfg := &FlagConfig{
		ConfigPath:      *cfgPath,
		SettingsDirPath: *settingsDir,
		ServicesDirPath: *servicesDir,
	}

	if err := validateFlags(flagCfg); err != nil {
		flag.Usage()
		return nil, err
	}

	return flagCfg, nil
}

// TODO: add suppressor flags
// --quiet
// --settings-only
// --installs-only
// --repos-only
// --commands-only
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
