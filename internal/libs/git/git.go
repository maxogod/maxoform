package git

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/libs/shell"
	"github.com/maxogod/maxoform/internal/utils"
)

func CloneMissingRepos(runner shell.Executor, repos []config.Repo, shouldSetupSSH bool) error {
	for _, repo := range repos {
		dest, err := utils.ExpandHome(repo.Dest)
		if err != nil {
			return fmt.Errorf("expanding %q: %w", repo.Dest, err)
		}

		if utils.PathExists(dest) {
			if err := runner.Run("git", "-C", dest, "pull", "--ff-only"); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating parent directory for %q: %w", dest, err)
		}

		if err := runner.Run("git", "clone", repo.URL, dest); err != nil {
			return err
		}

		if !shouldSetupSSH {
			continue
		}

		sshURL, err := toSSHRemote(repo.URL)
		if err != nil {
			return fmt.Errorf("converting remote to ssh for %q: %w", repo.URL, err)
		}
		if sshURL != "" {
			if err := runner.Run("git", "-C", dest, "remote", "set-url", "origin", sshURL); err != nil {
				return err
			}
		}
	}

	return nil
}

func ConfigureGlobalUser(runner shell.Executor, name, email, coreEditor string) error {
	if name != "" {
		if err := runner.Run("git", "config", "--global", "user.name", name); err != nil {
			return err
		}
	}

	if email != "" {
		if err := runner.Run("git", "config", "--global", "user.email", email); err != nil {
			return err
		}
	}

	if coreEditor != "" {
		if err := runner.Run("git", "config", "--global", "core.editor", coreEditor); err != nil {
			return err
		}
	}
	return nil
}

func toSSHRemote(raw string) (string, error) {
	if strings.HasPrefix(raw, "git@") {
		return "", nil
	}

	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return "", nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, ".git")
	if path == "" || u.Host == "" {
		return "", fmt.Errorf("invalid repository URL")
	}

	return fmt.Sprintf("git@%s:%s.git", u.Host, path), nil
}
