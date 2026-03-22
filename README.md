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

| Spell        | Incantation                              | Effect                                                   |
| ------------ | ---------------------------------------- | -------------------------------------------------------- |
| **Travel**   | `rift <name>`                            | Tear open a rift and teleport to a memorized waypoint    |
| **Discover** | `rift discover <name>`                   | Inscribe your current location as a waypoint             |
| **Awaken**   | `rift awaken`                            | Repair shell integration or database if something breaks |
| **Waypoint** | `rift waypoint [name] [--destroy\|--rebind\|--reforge]` | Inspect, rename, rebind, or destroy waypoints |

> Full usage, examples, and settings flags → **[usage.md](usage.md)**

---

## The Mechanics Behind the Magic

A child process cannot alter the working directory of its parent shell — this is a hard constraint of every OS. Rift works around it by inscribing a shell function into your config on first run. When you cast `rift <name>`, the function runs the binary and captures its output, then hands the `cd` command back to your live shell session via `eval`.

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

| Platform | Database                                        | Settings                                                         |
| -------- | ----------------------------------------------- | ---------------------------------------------------------------- |
| macOS    | `~/Library/Application Support/rift/db/rift.db` | `~/Library/Application Support/rift/settings/rift_settings.json` |
| Linux    | `~/.config/rift/db/rift.db`                     | `~/.config/rift/settings/rift_settings.json`                     |
| Windows  | `%APPDATA%\rift\db\rift.db`                     | `%APPDATA%\rift\settings\rift_settings.json`                     |

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

---

## Changelog

### v0.1.0

Core

- Discover — save a waypoint by giving a custom name to any directory path
- Navigate — jump to a saved waypoint instantly by name

Waypoint Management
- List — view all saved waypoints
- Info — view details of a specific waypoint
- Destroy — delete a waypoint
- Rebind — change the path of an existing waypoint
- Reforge — rename an existing waypoint

Settings Flag
- --update — check for updates
- --version — show current version
- --language — set the display language for rift
- --autoupdate — enable or disable automatic update checks on startup
- --download-pre-release — include pre-release versions when checking for updates
