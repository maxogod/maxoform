package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages PackagesConfig `yaml:"packages"`
	Repos    []Repo         `yaml:"repos"`
	Settings Settings       `yaml:"settings"`
}

type PackagesConfig struct {
	Apt  []string `yaml:"apt"`
	Snap []string `yaml:"snap"`
	Npm  []string `yaml:"npm"`
}

type Repo struct {
	URL  string `yaml:"url"`
	Dest string `yaml:"dest"`
}

type Settings struct {
	ImportDconf   bool   `yaml:"import_dconf"`
	GitUserName   string `yaml:"git_user_name"`
	GitUserEmail  string `yaml:"git_user_email"`
	SshKeyPath    string `yaml:"ssh_key_path"`
	SshPassphrase string
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

	passPhrase, ok := os.LookupEnv("MF_SSH_PATHPHRASE")
	if !ok {
		return nil, fmt.Errorf("MF_SSH_PATHPHRASE environment variable is not set")
	}
	cfg.Settings.SshPassphrase = passPhrase

	return &cfg, nil
}
