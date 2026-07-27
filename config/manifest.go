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
