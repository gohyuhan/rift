<p align="center">
  <img src="assets/rift.png" alt="rift logo" width="640">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/github/v/release/gohyuhan/rift?style=for-the-badge" alt="Release">
  <img src="https://img.shields.io/github/license/gohyuhan/rift?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/github/actions/workflow/status/gohyuhan/rift/release.yml?style=for-the-badge&label=Release" alt="Release">
</p>

<p align="center">
  <strong>Tear through your filesystem. No spells required.</strong>
</p>

---

## Lore

In the age of the terminal, brave developers were cursed to wander — chaining `cd` after `cd`, descending into folder after folder, losing precious time in the labyrinth of the filesystem.

Then came **rift**.

A tool forged for those who refuse to waste their power on traversal. Memorize a location once. Tear open a rift. Step through. The path bends to your will.

---

## Spellbook

| Spell        | Incantation              | Effect                                                 |
| ------------ | ------------------------ | ------------------------------------------------------ |
| **Travel**   | `rift <name>`            | Tear open a rift and teleport to a memorized location  |
| **Memorize** | `rift --memorize <name>` | Inscribe your current location into memory             |
| **Awaken**   | `rift awaken`            | Awaken rift within your shell — run once after install |

> More spells incoming — `rift atlas`, `rift forget`, `rift rebind`.

---

## Casting

### Travel — tear open a rift

```sh
rift <name>
```

```sh
# Lost deep in /Users/alice/projects/frontend/src/components/ui
rift api

# Rift tears open. You step through.
# You are now at /Users/alice/projects/backend/src/api
```

### Memorize — inscribe a location

Memorize your current directory and bind it to a name. Like a mage preparing spells before a long journey — do it once, call upon it forever.

```sh
rift --memorize <name>
```

```sh
cd /Users/alice/projects/backend/src/api
rift --memorize api
# Memorized: api -> /Users/alice/projects/backend/src/api
```

### Awaken — bind rift to your shell

Run once after installation. `rift awaken` detects your shell and inscribes the invocation function into its config file — the ritual that allows rift to manipulate your shell's path from within.

```sh
rift awaken
```

Supports **zsh**, **bash**, **fish**, **ksh**, and **PowerShell** across macOS, Linux, and Windows. After awakening, restart your shell or source the config.

---

## The Mechanics Behind the Magic

A child process cannot alter the working directory of its parent shell — this is a hard constraint of every OS. `rift awaken` works around it by inscribing a shell function into your config. When you cast `rift <name>`, the function runs the binary and captures its output, then hands the `cd` command back to your live shell session via `eval`.

```
You cast: rift api
         │
         ▼
shell function `rift` intercepts
         │  invokes the rift binary, captures stdout
         ▼
rift binary resolves the memorized path
         │  emits to stdout:  cd "/Users/alice/projects/backend/src/api"
         ▼
shell function evals the output
         │
         ▼
your shell changes directory — the rift closes behind you
```

Everything else — errors, confirmations, help text — travels through **stderr** and is displayed directly, never evaluated. Only one thing ever reaches stdout: the `cd` command itself.

```sh
# zsh / bash
rift() { eval "$(command rift "$@")"; }

# fish
function rift
    eval (command rift $argv)
end

# PowerShell (Windows)
function rift { Invoke-Expression (rift.exe $args) }

# PowerShell (macOS / Linux)
function rift { Invoke-Expression (& (Get-Command -CommandType Application rift) $args) }
```

See [`docs/shell-integration.md`](docs/shell-integration.md) for the full arcane breakdown.

---

## Grimoire

Memorized locations are stored in a [bbolt](https://github.com/etcd-io/bbolt) database — your personal grimoire, persisted across sessions.

| Platform | Database | Settings |
| -------- | -------- | -------- |
| macOS | `~/Library/Application Support/rift/db/rift.db` | `~/Library/Application Support/rift/settings/rift_settings.json` |
| Linux | `~/.config/rift/db/rift.db` | `~/.config/rift/settings/rift_settings.json` |
| Windows | `%APPDATA%\rift\db\rift.db` | `%APPDATA%\rift\settings\rift_settings.json` |

## Installation

### Linux

```bash
curl --proto "=https" -sSfL https://github.com/gohyuhan/rift/releases/latest/download/install.sh | bash
```

### macOS (curl or homebrew)

```bash
curl --proto "=https" -sSfL https://github.com/gohyuhan/rift/releases/latest/download/install.sh | bash

# via homebrew
# Add the tap (once)
brew tap gohyuhan/rift

# Install latest
brew update && brew install rift
```

### Windows (PowerShell or scoop)

```powershell
powershell -c "irm https://github.com/gohyuhan/rift/releases/latest/download/install.ps1 | iex"

# via scoop
# Add the bucket (once)
scoop bucket add rift https://github.com/gohyuhan/scoop-rift

# Install latest
scoop update; scoop install rift
```

### Go Install

If you have Go installed, you can install rift directly:

```bash
go install github.com/gohyuhan/rift@latest
```

## Uninstall & Cleanup

### macOS (Homebrew)

```bash
# 1. Uninstall + remove ALL versions
brew uninstall --force rift

# 2. Remove the tap
brew untap gohyuhan/rift

# 3. Delete the binary directly (in case it's not a symlink or brew missed it)
rm -f /opt/homebrew/bin/rift
rm -f /usr/local/bin/rift

# 4. Delete the entire Cellar folder for rift (old kegs)
rm -rf /opt/homebrew/Cellar/rift
rm -rf /usr/local/Cellar/rift

# 5. Delete any leftover symlinks
rm -rf /opt/homebrew/opt/rift
rm -rf /usr/local/opt/rift

# 6. Delete all cached downloads for rift
rm -rf ~/Library/Caches/Homebrew/rift*
rm -rf ~/Library/Caches/Homebrew/downloads/*rift*
```

### Windows (Scoop)

```powershell
# 1. Uninstall the app (all versions)
scoop uninstall rift 2>$null

# 2. Remove the bucket
scoop bucket rm rift 2>$null

# 3. Delete the app folder completely (including shims + persist)
rm -r -force "$env:USERPROFILE\scoop\apps\rift" 2>$null

# 4. Delete the bucket clone
rm -r -force "$env:USERPROFILE\scoop\buckets\rift" 2>$null

# 5. Delete all cached installers for rift
scoop cache rm "rift*" 2>$null
```

### Manual Installation (curl / powershell)

#### macOS / Linux

```bash
# Remove binary (if installed via curl)
sudo rm -f /usr/local/bin/rift
```

#### Windows

```powershell
# Remove binary and directory
Remove-Item -Path "$env:LOCALAPPDATA\rift" -Recurse -Force
```

### Configuration Cleanup

To completely remove rift's configuration files:

#### macOS

```bash
rm -rf "$HOME/Library/Application Support/rift"
```

#### Linux

```bash
rm -rf "$HOME/.config/rift"
```

#### Windows

```powershell
Remove-Item -Path "$env:APPDATA\rift" -Recurse -Force
```
