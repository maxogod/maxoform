package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
	"github.com/maxogod/maxoform/internal/utils"
)

func InstallAndEnable(runner shell.Executor, servicesDir string, manifest *config.ServiceManifest) error {
	for idx, serviceFile := range manifest.Entries {
		if strings.TrimSpace(serviceFile) == "" {
			return fmt.Errorf("service manifest entry at index %d is empty", idx)
		}

		sourcePath := filepath.Join(servicesDir, serviceFile)
		if !utils.PathExists(sourcePath) {
			return fmt.Errorf("service file not found: %s", sourcePath)
		}

		targetPath := filepath.Join("/etc/systemd/system", serviceFile)
		if err := runner.Run("sudo", "cp", sourcePath, targetPath); err != nil {
			return err
		}
	}

	if len(manifest.Entries) == 0 {
		logger.Log.Info("no services configured in manifest, skipping")
		return nil
	}

	if err := runner.Run("sudo", "systemctl", "daemon-reload"); err != nil {
		return err
	}

	for _, serviceFile := range manifest.Entries {
		if err := runner.Run("sudo", "systemctl", "enable", "--now", serviceFile); err != nil {
			return err
		}
	}

	return nil
}
