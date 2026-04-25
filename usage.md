# rift — Usage Guide

Everything you need to wield rift from the terminal.

---

## Table of Contents

- [First-time Setup](#first-time-setup)
- [Commands](#commands)
  - [awaken](#awaken)
  - [discover](#discover)
  - [travel](#travel)
  - [waypoint](#waypoint)
  - [learn](#learn)
  - [spell](#spell)
  - [rune](#rune)
  - [spellbook](#spellbook)
- [Settings Flags](#settings-flags)
  - [--language](#--language)
  - [--autoupdate](#--autoupdate)
  - [--download-pre-release](#--download-pre-release)
  - [--update](#--update)
  - [--version](#--version)
- [Supported Languages](#supported-languages)

---

## First-time Setup

Rift sets itself up automatically on first run — shell integration and database initialization happen in the background without any extra step from you. Just install and start using it.

---

## Commands

### awaken

Repairs rift's shell integration and database. Rift sets up automatically on first run, so you rarely need this — but if something goes wrong (waypoints not resolving, shell navigation broken, database errors), running `rift awaken` will fix it in most cases.

```sh
rift awaken
```

**When to use it**
- Shell navigation stops working (the `rift <name>` command no longer changes directory)
- You see database-related errors on startup
- You switched to a new shell and rift is not recognized in it
- You reinstalled or moved the rift binary

Supported shells: **zsh**, **bash**, **fish**, **ksh**, **PowerShell**
Supported platforms: macOS, Linux, Windows

---

### discover

Assigns a waypoint name to your current working directory and saves it for fast navigation.

```sh
rift discover <name>
```

**Examples**

```sh
# Navigate to the directory you want to save, then discover it
cd /Users/alice/projects/backend/src/api
rift discover api
# rift: saved "api" -> /Users/alice/projects/backend/src/api

cd /Users/alice/projects/frontend/src/components
rift discover ui
# rift: saved "ui" -> /Users/alice/projects/frontend/src/components
```

**Notes**
- Waypoint names are unique. Attempting to discover a name that already exists will return an error.
- Reserved keywords (`rift`, `awaken`, `discover`, `waypoint`, `learn`, `spell`, and others) cannot be used as waypoint names.

---

### travel

Tears open a rift and teleports you to the directory bound to a waypoint name. Pass `--cast` to run a spell at the waypoint's path instead of navigating there.

```sh
rift <name>
rift <name> --cast <spell-name-or-command>
```

**Flags**

| Flag | Description |
| ---- | ----------- |
| `--cast <spell-name-or-command>` | Instead of navigating, cast a learned spell or run a command string at the waypoint's path as the working directory |

**Examples**

```sh
# You are somewhere deep in the filesystem
rift api
# You are now at /Users/alice/projects/backend/src/api

rift ui
# You are now at /Users/alice/projects/frontend/src/components

# Cast a saved spell at a waypoint's path without changing your current directory
rift api --cast build
# runs: docker compose up --build (in /Users/alice/projects/backend/src/api)

# Run a command string directly at a waypoint's path (no saved spell needed)
rift api --cast "git status"
# runs: git status (in /Users/alice/projects/backend/src/api)
```

**Notes**
- If the waypoint name does not exist, rift will tell you.
- Waypoints persist across sessions — discover once, travel forever.
- `--cast` first tries to resolve the value as a saved spell name; if no match is found, it parses and runs the value as a raw command string.
- When `--cast` is used, your working directory does not change.

---

### waypoint

Inspects and manages your saved waypoints.

```sh
rift waypoint
rift waypoint <name>
rift waypoint <name> --destroy
rift waypoint <name> --rebind
rift waypoint <name> --rebind=<path>
rift waypoint <name> --reforge <new-name>
```

With no arguments, lists every saved waypoint with its name and bound path. With a waypoint name, shows detailed info for that waypoint. Flags unlock destructive or mutating operations.

**Flags**

| Flag | Description |
| ---- | ----------- |
| `--destroy` | Permanently removes the named waypoint |
| `--rebind` | Reassigns the waypoint to the current working directory |
| `--rebind=<path>` | Reassigns the waypoint to the given absolute path |
| `--reforge <new-name>` | Renames the waypoint to a new name, preserving all its data |

**Examples**

```sh
# List all waypoints
rift waypoint

# Show detail for a specific waypoint
rift waypoint api

# Destroy a waypoint
rift waypoint api --destroy
# rift: waypoint "api" has been destroyed

# Rebind a waypoint to the current directory
rift waypoint api --rebind
# rift: waypoint "api" rebound to /Users/alice/projects/new-api

# Rebind a waypoint to an explicit path
rift waypoint api --rebind=/Users/alice/projects/other-api
# rift: waypoint "api" rebound to /Users/alice/projects/other-api

# Rename a waypoint
rift waypoint api --reforge backend-api
# rift: waypoint "api" reforged to "backend-api"
```

**Notes**
- `--destroy`, `--rebind`, and `--reforge` are mutually exclusive — only one may be used at a time.
- `--rebind` resets the waypoint's sealed state, sealed reason, and travelled count, and updates its discovered timestamp.
- `--reforge` preserves all waypoint data (path, sealed state, travelled count, timestamps) — only the name changes.
- `--reforge` requires a non-empty new name and will refuse to overwrite an existing waypoint.

#### Interactive TUI

Running `rift waypoint` with no arguments launches an interactive browser for all saved waypoints.

**Key bindings**

| Key | Action |
| --- | ------ |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `enter` | Travel to selected waypoint |
| `r` | Open rebind path input |
| `R` | Open reforge name input |
| `u` / `U` | Unseal selected waypoint |
| `y` | Copy waypoint name to clipboard |
| `Y` | Copy waypoint path to clipboard |
| `backspace` | Destroy selected waypoint |
| `?` | Open help popup |

In rebind/reforge input popups:

| Key | Action |
| --- | ------ |
| `ctrl+y` | Copy input field text |
| `ctrl+p` | Paste from clipboard |

---

### learn

Binds a terminal command to a spell name and saves it for quick recall. If the spell name already exists, its command is overwritten and its cast count is reset.

```sh
rift learn <spell name> <command>
```

Multi-word commands must be wrapped in quotes.

**Examples**

```sh
rift learn build "docker compose up --build"
# rift: learned "build" -> docker compose up --build

rift learn test "go test ./..."
# rift: learned "test" -> go test ./...

# Override an existing spell
rift learn build "make build"
# rift: spell "build" updated -> make build
```

**Notes**
- Spell names are subject to the same reserved keyword restrictions as waypoint names.
- Spell names cannot contain whitespace.
- The spell command cannot contain `cd` — use `discover` + `rift <name>` for navigation instead.
- Overriding an existing spell resets its cast count to 0.

---

### spell

Casts a learned spell by name, or runs a command string directly in the current working directory. Pass `--forget` to remove a saved spell instead.

```sh
rift spell <name-or-command>
rift spell <name> --forget
```

**Flags**

| Flag | Description |
| ---- | ----------- |
| `--forget` | Permanently removes the named spell from rift's memory |

**Examples**

```sh
# Cast a saved spell
rift spell build
# runs: docker compose up --build

rift spell test
# runs: go test ./...

# Run a command string directly (no saved spell needed)
rift spell "git log --oneline -10"
# runs: git log --oneline -10

# Remove a spell
rift spell build --forget
# rift: spell "build" forgotten
```

**Notes**
- If the value matches a saved spell name, that spell is cast (and its cast count incremented). Otherwise, rift parses and runs the value as a raw command string.
- The command runs in the directory you are in when you invoke `rift spell` — it does not change your working directory.
- The exit code of the bound command is not propagated; rift is a launcher, not a validator of the command's outcome.
- `--forget` only works with saved spell names, not raw command strings.
- `--forget` is permanent — the spell cannot be recovered after removal.

---

### rune

Engraves on-enter and on-leave shell commands onto a waypoint. When you travel **to** the waypoint via `rift <name>`, its on-enter commands run automatically. When you travel **away** from a waypoint (your current directory matches its path and you `rift` elsewhere), its on-leave commands run first. Runes only fire when traveling via rift — plain `cd` does not trigger them.

```sh
rift rune <waypoint-name>
```

**Examples**

```sh
# Engrave rune commands onto the "api" waypoint
rift rune api

# After engraving on-enter, traveling fires the rune automatically
rift api
# [RUNE_ON_ENTER (1/2) - docker compose up -d]
# ...
# [RUNE_ON_ENTER (2/2) - echo "env loaded"]
# ...
```

**Notes**
- Rune commands cannot use `cd` — path changes inside runes have no effect. Use `rift <name>` for navigation instead.
- Commands run with the waypoint's path as the working directory.
- Multiple commands per slot — enter one command per line in the textarea.
- On-leave fires based on your current working directory at travel time, not the waypoint name.
- Engraving a slot overwrites any previously saved commands for that slot.

#### Interactive TUI

Running `rift rune <waypoint-name>` opens an interactive TUI for engraving and removing rune commands.

**Option selection**

| Key | Action |
| --- | ------ |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `enter` | Select option |
| `q` / `ctrl+c` | Quit |
| `esc` | Quit |

**Available options**

| Option | Description |
| ------ | ----------- |
| Engrave On-Enter Rune | Set commands to run when entering this waypoint |
| Engrave On-Leave Rune | Set commands to run when leaving this waypoint |
| Remove On-Enter Rune | Clear all on-enter commands from this waypoint |
| Remove On-Leave Rune | Clear all on-leave commands from this waypoint |

**In command entry popup**

| Key | Action |
| --- | ------ |
| `tab` | Move focus from textarea to Engrave button |
| `shift+tab` | Move focus back to textarea |
| `ctrl+y` | Copy textarea content to clipboard |
| `ctrl+p` | Paste from clipboard into textarea |
| `enter` | Save the rune (when Engrave button is focused) |
| `esc` | Go back to option list |

---

### spellbook

Browse and manage your learned spells through an interactive TUI. With a spell name, shows detailed info for that spell instead.

```sh
rift spellbook
rift spellbook <spell-name>
```

**Examples**

```sh
# Launch the interactive spellbook
rift spellbook

# Show detail for a specific spell
rift spellbook build
```

#### Interactive TUI

Running `rift spellbook` with no arguments launches an interactive browser for all learned spells.

**Key bindings**

| Key | Action |
| --- | ------ |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `enter` | Cast the selected spell (opens location chooser) |
| `n` / `N` | Open learn popup to add a new spell |
| `backspace` | Forget the selected spell |
| `?` | Open help popup |
| `esc` | Close active popup, or quit if no popup is open |
| `q` / `ctrl+c` | Quit |

In the cast location popup:

| Key | Action |
| --- | ------ |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `enter` | Confirm selected cast location |
| `esc` | Close popup |

In the learn spell popup:

| Key | Action |
| --- | ------ |
| `tab` | Move to next input field |
| `shift+tab` | Move to previous input field |
| `enter` | Save the new spell |
| `ctrl+y` | Copy current input field value |
| `ctrl+p` | Paste into current input field |
| `esc` | Close popup without saving |

**Notes**
- Casting from the spellbook lets you choose between the current working directory or any discovered waypoint as the execution path.
- Your working directory does not change when casting from the spellbook.
- Forgetting a spell via `backspace` is permanent.

---

## Settings Flags

Settings flags are passed to the root `rift` command (without a waypoint name). They update your rift configuration and are persisted across sessions.

```sh
rift --<flag>
rift --<flag>=<value>
```

---

### --language

Sets the display language for all rift output.

```sh
rift --language <code>
```

| Code | Language |
| ---- | -------- |
| `EN` | English |
| `JA` | Japanese |
| `ZH-HANS` | Simplified Chinese |
| `ZH-HANT` | Traditional Chinese |

**Examples**

```sh
rift --language JA
# rift: 言語を JA に設定しました

rift --language ZH-HANS
# rift: 语言已设置为 ZH-HANS
```

The language change takes effect immediately and applies to all subsequent rift output.

---

### --autoupdate

Controls whether rift automatically checks for new releases on startup.

```sh
# Enable (default)
rift --autoupdate

# Disable
rift --autoupdate=false
```

When enabled, rift checks for a newer release each time it runs (respecting a cooldown interval). If a newer version is found, you will be prompted to download it.

**Examples**

```sh
rift --autoupdate
# rift: auto-update set to true

rift --autoupdate=false
# rift: auto-update set to false
```

---

### --download-pre-release

Controls whether rift considers pre-release versions when checking for updates.

```sh
# Enable — rift will offer pre-release versions
rift --download-pre-release

# Disable (default) — stable releases only
rift --download-pre-release=false
```

This flag only has an effect when `--autoupdate` is enabled.

**Examples**

```sh
rift --download-pre-release
# rift: download pre-release set to true

rift --download-pre-release=false
# rift: download pre-release set to false
```

---

### --update

Manually triggers an update check and downloads the latest version if one is available. Unlike `--autoupdate`, this bypasses the cooldown interval and runs immediately.

```sh
rift --update
```

**Examples**

```sh
# Already on the latest version
rift --update
# You are already on the latest version (v0.1.0)

# A newer version is available
rift --update
# Downloading version v0.2.0...
# Successfully updated to version v0.2.0
```

**Notes**
- Does not prompt for confirmation — update proceeds immediately if a newer version is found.
- Respects the `--download-pre-release` setting when deciding which release to compare against.

---

### --version

Prints the current installed version of rift and exits.

```sh
rift --version
```

**Examples**

```sh
rift --version
# v0.1.0
```

---

## Supported Languages

| Code | Language |
| ---- | -------- |
| `EN` | English (default) |
| `JA` | Japanese |
| `ZH-HANS` | Simplified Chinese |
| `ZH-HANT` | Traditional Chinese |

To change the language:

```sh
rift --language <code>
```
