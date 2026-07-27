package sshkey

import (
	"fmt"
	"os"
	"strings"

	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/utils"
)

func EnsureAndPrint(runner shell.Runner, email, keyPath, passPhrase string) error {
	if email == "" {
		return fmt.Errorf("ssh key email is required")
	}

	expandedKeyPath, err := utils.ExpandHome(keyPath)
	if err != nil {
		return fmt.Errorf("resolving public key path: %w", err)
	}
	publicKeyPath := expandedKeyPath + ".pub"

	exists := utils.PathExists(publicKeyPath)
	if !exists {
		err := runner.Run(
			"ssh-keygen",
			"-t", "ed25519",
			"-C", email,
			"-f", expandedKeyPath,
			"-N", passPhrase,
		)
		if err != nil {
			return err
		}
	}

	pub, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("reading public key %s: %w", publicKeyPath, err)
	}
	fmt.Println(strings.TrimSpace(string(pub)))
	return nil
}
