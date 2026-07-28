package dconf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/shell"
)

func Import(runner shell.Runner, settingsDir string, manifest *config.DconfManifest) error {
	for _, e := range manifest.Entries {
		if err := loadFile(runner, e.Key, filepath.Join(settingsDir, e.File)); err != nil {
			return err
		}
	}

	return nil
}

func loadFile(runner shell.Runner, key, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening %s: %w", filePath, err)
	}
	defer f.Close()

	if err := runner.RunWithStdin(f, "dconf", "load", key); err != nil {
		return fmt.Errorf("loading dconf key %s from %s: %w", key, filePath, err)
	}

	return nil
}
