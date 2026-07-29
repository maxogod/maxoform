package config

import (
	"flag"
	"fmt"

	"github.com/maxogod/maxoform/internal/utils"
)

type FlagConfig struct {
	ConfigPath        string
	SettingsDirPath   string
	DconfManifestPath string
}

func LoadFlagConf() (*FlagConfig, error) {
	cfgPath := flag.String("config", "", "path to packages YAML")
	settingsDir := flag.String("settings-dir", "", "path to dconf .ini files directory")
	manifestPath := flag.String("dconf-manifest", "", "path to dconf manifest YAML")
	flag.Parse()

	flagCfg := &FlagConfig{
		ConfigPath:        *cfgPath,
		SettingsDirPath:   *settingsDir,
		DconfManifestPath: *manifestPath,
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
	if flagCfg.DconfManifestPath == "" || !utils.PathExists(flagCfg.DconfManifestPath) {
		return fmt.Errorf("missing --dconf-manifest")
	}
	return nil
}
