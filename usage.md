# rift — Usage Guide

Everything you need to wield rift from the terminal.

---

## Table of Contents

- [First-time Setup](#first-time-setup)
- [Commands](#commands)
  - [awaken](#awaken)
  - [discover](#discover)
  - [travel](#travel)
- [Settings Flags](#settings-flags)
  - [--language](#--language)
  - [--autoupdate](#--autoupdate)
  - [--download-pre-release](#--download-pre-release)
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
- Reserved keywords (`rift`, `awaken`, `discover`) cannot be used as waypoint names.

---

### travel

Tears open a rift and teleports you to the directory bound to a waypoint name.

```sh
rift <name>
```

**Examples**

```sh
# You are somewhere deep in the filesystem
rift api
# You are now at /Users/alice/projects/backend/src/api

rift ui
# You are now at /Users/alice/projects/frontend/src/components
```

**Notes**
- If the waypoint name does not exist, rift will tell you.
- Waypoints persist across sessions — discover once, travel forever.

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
