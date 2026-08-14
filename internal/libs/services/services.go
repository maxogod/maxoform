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
	servicesToInstall := make([]string, 0, len(manifest.Entries))

	for idx, entry := range manifest.Entries {
		serviceFile := strings.TrimSpace(entry.Service)
		condition := strings.TrimSpace(entry.Condition)

		if serviceFile == "" {
			return fmt.Errorf("service manifest entry at index %d has empty service", idx)
		}
		if condition == "" {
			return fmt.Errorf("service manifest entry at index %d has empty condition", idx)
		}

		if !runner.CheckShell(condition) {
			logger.Log.Infof("service condition not met, skipping %s: %s", serviceFile, condition)
			continue
		}

		sourcePath := filepath.Join(servicesDir, serviceFile)
		if !utils.PathExists(sourcePath) {
			return fmt.Errorf("service file not found: %s", sourcePath)
		}

		targetPath := filepath.Join("/etc/systemd/system", serviceFile)
		if err := runner.Run("sudo", "cp", sourcePath, targetPath); err != nil {
			return err
		}

		servicesToInstall = append(servicesToInstall, serviceFile)
	}

	if len(servicesToInstall) == 0 {
		logger.Log.Info("no services to install after evaluating manifest conditions, skipping")
		return nil
	}

	if err := runner.Run("sudo", "systemctl", "daemon-reload"); err != nil {
		return err
	}

	for _, serviceFile := range servicesToInstall {
		if err := runner.Run("sudo", "systemctl", "enable", "--now", serviceFile); err != nil {
			return err
		}
	}

	return nil
}
