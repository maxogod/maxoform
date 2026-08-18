package swap

import (
	"fmt"

	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/logger"
)

const filePath = "/swapfile"

// Ensure creates and enables a /swapfile of the given size, persisting it in
// /etc/fstab so it survives reboots. If sizeGB is 0 or less, swap setup is
// skipped entirely. If a swap file of the requested size is already active,
// nothing is done.
func Ensure(runner shell.Executor, sizeGB int) error {
	if sizeGB <= 0 {
		logger.Log.Info("swap size not configured, skipping swap file setup")
		return nil
	}

	sizeBytes := int64(sizeGB) * 1024 * 1024 * 1024
	activeCheck := fmt.Sprintf("swapon --show=NAME --noheadings | grep -qx %s", filePath)
	sizeCheck := fmt.Sprintf("[ \"$(stat -c%%s %s 2>/dev/null)\" = \"%d\" ]", filePath, sizeBytes)

	active := runner.CheckShell(activeCheck)
	if active && runner.CheckShell(sizeCheck) {
		logger.Log.Infof("swap file already active at %dG, skipping: %s", sizeGB, filePath)
		return nil
	}

	if active {
		if err := runner.Run("sudo", "swapoff", filePath); err != nil {
			return fmt.Errorf("disabling existing swap file: %w", err)
		}
	}

	if err := runner.Run("sudo", "fallocate", "-l", fmt.Sprintf("%dG", sizeGB), filePath); err != nil {
		return fmt.Errorf("allocating swap file: %w", err)
	}
	if err := runner.Run("sudo", "chmod", "600", filePath); err != nil {
		return fmt.Errorf("setting swap file permissions: %w", err)
	}
	if err := runner.Run("sudo", "mkswap", filePath); err != nil {
		return fmt.Errorf("formatting swap file: %w", err)
	}
	if err := runner.Run("sudo", "swapon", filePath); err != nil {
		return fmt.Errorf("enabling swap file: %w", err)
	}

	fstabCheck := fmt.Sprintf("grep -qF %s /etc/fstab", filePath)
	if !runner.CheckShell(fstabCheck) {
		fstabEntry := fmt.Sprintf("echo '%s none swap sw 0 0' | sudo tee -a /etc/fstab >/dev/null", filePath)
		if err := runner.RunShell(fstabEntry); err != nil {
			return fmt.Errorf("persisting swap file in /etc/fstab: %w", err)
		}
	}

	return nil
}
