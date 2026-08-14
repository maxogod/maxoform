package config

import (
	"flag"
	"fmt"

	"github.com/maxogod/maxoform/internal/utils"
)

type FlagConfig struct {
	ConfigPath           string
	SettingsDirPath      string
	DconfManifestPath    string
	ServicesDirPath      string
	ServicesManifestPath string
}

func LoadFlagConf() (*FlagConfig, error) {
	cfgPath := flag.String("config", "", "path to packages YAML")
	settingsDir := flag.String("settings-dir", "", "path to dconf .ini files directory")
	manifestPath := flag.String("dconf-manifest", "", "path to dconf manifest YAML")
	servicesDir := flag.String("services-dir", "data/services", "path to service files directory")
	servicesManifest := flag.String("services-manifest", "data/services/manifest.yaml", "path to services manifest YAML")
	flag.Parse()

	flagCfg := &FlagConfig{
		ConfigPath:           *cfgPath,
		SettingsDirPath:      *settingsDir,
		DconfManifestPath:    *manifestPath,
		ServicesDirPath:      *servicesDir,
		ServicesManifestPath: *servicesManifest,
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
	if flagCfg.ServicesDirPath == "" || !utils.PathExists(flagCfg.ServicesDirPath) {
		return fmt.Errorf("missing --services-dir")
	}
	if flagCfg.ServicesManifestPath == "" || !utils.PathExists(flagCfg.ServicesManifestPath) {
		return fmt.Errorf("missing --services-manifest")
	}
	return nil
}
