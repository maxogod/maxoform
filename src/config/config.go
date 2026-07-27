package config

import (
	"gopkg.in/yaml.v3"
	"os"
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
	ImportDconf bool `yaml:"import_dconf"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
