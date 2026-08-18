package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages     PackagesConfig     `yaml:"packages"`
	Repos        []Repo             `yaml:"repos"`
	Settings     Settings           `yaml:"settings"`
	NpmBootstrap NpmBootstrapConfig `yaml:"npm_bootstrap"`
	Commands     CommandsConfig     `yaml:"commands"`
}

// TODO: add a pre-installation commands list to execute before apt/snap/npm install
type PackagesConfig struct {
	Apt  []string      `yaml:"apt"`
	Snap []SnapPackage `yaml:"snap"`
	Npm  []string      `yaml:"npm"`
	Pipx []string      `yaml:"pipx"`
}

type SnapPackage struct {
	Name    string `yaml:"name"`
	Classic bool   `yaml:"classic"`
}

type Repo struct {
	URL  string `yaml:"url"`
	Dest string `yaml:"dest"`
}

type NpmBootstrapConfig struct {
	Enabled          bool   `yaml:"enabled"`
	InstallScriptURL string `yaml:"install_script_url"`
	NodeVersion      string `yaml:"node_version"`
}

type CommandsConfig struct {
	Post []ShellCommand `yaml:"post"`
}

type ShellCommand struct {
	Run string `yaml:"run"`
}

type Settings struct {
	ImportDconf   bool   `yaml:"import_dconf"`
	GitUserName   string `yaml:"git_user_name"`
	GitUserEmail  string `yaml:"git_user_email"`
	GitCoreEditor string `yaml:"git_core_editor"`
	SshKeyPath    string `yaml:"ssh_key_path"`
	SshPassphrase string
	SwapFilePath  string `yaml:"swap_file_path"`
	SwapSizeGB    int    `yaml:"swap_size_gb"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	passPhrase, ok := os.LookupEnv("MF_SSH_PASSPHRASE")
	if !ok {
		return nil, fmt.Errorf("MF_SSH_PASSPHRASE environment variable is not set")
	}
	cfg.Settings.SshPassphrase = passPhrase

	return &cfg, nil
}
