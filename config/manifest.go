package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DconfManifest struct {
	Entries []DconfEntry `yaml:"entries"`
}

type DconfEntry struct {
	Key  string `yaml:"key"`
	File string `yaml:"file"`
}

type ServiceManifest struct {
	Entries []ServiceEntry `yaml:"entries"`
}

type ServiceEntry struct {
	Condition string `yaml:"condition"`
	Service   string `yaml:"service"`
}

func LoadDconfManifest(path string) (*DconfManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest DconfManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func LoadServiceManifest(path string) (*ServiceManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest ServiceManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}
