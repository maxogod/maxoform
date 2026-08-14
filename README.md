# maxoform

`maxoform` is a Go-based Ubuntu setup tool that:

- installs packages from **apt**, **snap**, **npm**, and **pipx**
- clones Git repositories into configured destinations
- imports GNOME dconf settings from dumped `.ini` files
- configures global git identity and outputs your SSH public key (for copy-paste)
- runs ordered post-install shell commands
- installs and enables configured `systemd` service units

The tool is driven by data already in this repository:

- package/repo/commands config: `data/configuration/packages.yaml`
- dconf import manifest: `data/settings/manifest.yaml`
- dconf settings: `data/settings/*.ini`
- service manifest: `data/services/manifest.yaml`
- service unit files: `data/services/*.service`
- reference scripts for dconf dump/load: `scripts/dconf-dump-load.sh`

## How to use

Select a version and install and run it:

```bash
wget https://github.com/maxogod/maxoform/releases/download/v1.0.0/maxoform-linux-amd64.tar.gz
tar xzf maxoform-linux-amd64.tar.gz
cd maxoform-linux-amd64
MF_SSH_PASSPHRASE=mypassphrase ./maxoform --config data/configuration/packages.yaml --settings-dir data/settings --services-dir data/services
```

Flags:

- `--config`: path to YAML configuration file
- `--settings-dir`: directory containing dconf `.ini` files
- `--services-dir`: directory containing service unit files

Manifests are inferred automatically as:
- `<settings-dir>/manifest.yaml`
- `<services-dir>/manifest.yaml`

## Configuration format

`data/configuration/packages.yaml`:

```yaml
packages:
  apt: []
  snap: []
  npm: []
  pipx: []

npm_bootstrap:
  enabled: false
  install_script_url: https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.6/install.sh
  node_version: "26"

commands:
  post:
    - run: sudo snap remove firefox
    - run: sudo add-apt-repository -y ppa:obsproject/obs-studio
    - run: sudo apt install -y obs-studio

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
2. installs configured apt/snap/npm/pipx packages
3. clones missing repos only
4. imports dconf keys from `settings/manifest.yaml`
5. ensures `~/.ssh/id_ed25519.pub` exists and prints the key content
6. runs `commands.post` shell commands in order
7. copies and enables services from `data/services/manifest.yaml`


## Notes

- `apt` and `snap` commands are run with `sudo`.
- dconf import assumes GNOME keys/apps exist on the target Ubuntu system.
- SSH passphrase is read from `MF_SSH_PASSPHRASE`.
