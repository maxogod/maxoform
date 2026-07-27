# maxoform

`maxoform` is a Go-based Ubuntu setup tool that:

- installs packages from **apt**, **snap**, and **npm**
- clones Git repositories into configured destinations
- imports GNOME dconf settings from dumped `.ini` files
- configures global git identity and outputs your SSH public key (for copy-paste)

The tool is driven by data already in this repository with my settings:

- package/repo config: `data/config/packages.yaml`
- dconf import manifest: `data/settings/manifest.yaml`
- dconf settings: `data/settings/*.ini`
- reference scripts for dconf dump/load: `scripts/dconf-dump-load.sh`

## How to use

From the repository root:

```bash
go run ./src/cmd --config data/config/packages.yaml --settings-dir data/settings --dconf-manifest data/settings/manifest.yaml
```

Flags:

- `--config`: path to YAML configuration file
- `--settings-dir`: directory containing dconf `.ini` files
- `--dconf-manifest`: path to dconf manifest YAML (`entries[].key` + `entries[].file`)

## Configuration format

`config/packages.yaml`:

```yaml
packages:
  apt: []
  snap: []
  npm: []

repos:
  - url: git@github.com:example/repo.git
    dest: ~/.config/example

settings:
  import_dconf: true
  git_user_name: "example"
  git_user_email: "example@example.com"
```

## Implementation details

The program runs in phases:

1. updates system package state (`apt update/upgrade/autoremove`, `snap refresh`)
2. installs configured apt/snap/npm packages
3. clones missing repos only
4. imports dconf keys from `settings/manifest.yaml`
5. ensures `~/.ssh/id_ed25519.pub` exists and prints the key content

## Notes

- `apt` and `snap` commands are run with `sudo`.
- dconf import assumes GNOME keys/apps exist on the target Ubuntu system.
