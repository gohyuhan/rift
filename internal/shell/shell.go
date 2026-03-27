package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/logger"
	"github.com/gohyuhan/rift/style"
)

type Shell string

const (
	Zsh         Shell = "zsh"
	Bash        Shell = "bash"
	Fish        Shell = "fish"
	Ksh         Shell = "ksh"
	PowerShell  Shell = "powershell"
	CMD         Shell = "cmd" // Windows Command Prompt — not supported
	Unsupported Shell = "unsupported"
)

// ----------------------------------
//
//	Returns the current shell by inspecting environment variables in priority order:
//	  1. $SHELL suffix      — set on macOS, Linux, WSL, Git Bash, Cygwin, MSYS2
//	  2. $PSModulePath      — always present inside any PowerShell session (all platforms)
//	  3. $COMSPEC           — Windows CMD (only when PSModulePath is absent)
//	  4. OS-level fallback  — Linux/macOS defaults to bash; unknown Windows terminals → Unsupported
//
// ----------------------------------
func Detect() Shell {
	s := os.Getenv("SHELL")
	switch {
	case strings.HasSuffix(s, "zsh"):
		return Zsh
	case strings.HasSuffix(s, "bash"):
		return Bash
	case strings.HasSuffix(s, "fish"):
		return Fish
	case strings.HasSuffix(s, "ksh"), strings.HasSuffix(s, "mksh"):
		return Ksh
	case strings.Contains(s, "pwsh"), strings.Contains(s, "powershell"):
		return PowerShell
	case strings.HasSuffix(s, "csh"), strings.HasSuffix(s, "tcsh"):
		return Unsupported
	}

	// $SHELL is not set on native Windows terminals.
	if runtime.GOOS == "windows" {
		// PSModulePath is injected by PowerShell into every PS session.
		if os.Getenv("PSModulePath") != "" {
			return PowerShell
		}
		// COMSPEC points to cmd.exe — we are inside Command Prompt.
		if os.Getenv("COMSPEC") != "" {
			return CMD
		}
		// Unknown Windows terminal (e.g. Nushell, Elvish).
		return Unsupported
	}

	// Linux / macOS: $SHELL unset is unusual but fall back to bash.
	return Bash
}

// ----------------------------------
//
//	Returns the shell config file path for the given shell type.
//	Returns an error for CMD and Unsupported shells, which have no config file.
//
// ----------------------------------
func ConfigFile(sh Shell) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch sh {
	case Zsh:
		return filepath.Join(home, ".zshrc"), nil

	case Bash:
		// macOS login shells source .bash_profile; Linux interactive shells source .bashrc.
		// Git Bash / Cygwin / MSYS2 on Windows also use .bashrc (GOOS=windows but $SHELL set).
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".bash_profile"), nil
		}
		return filepath.Join(home, ".bashrc"), nil

	case Ksh:
		// ksh reads $ENV at startup; ~/.kshrc is the conventional interactive config.
		return filepath.Join(home, ".kshrc"), nil

	case Fish:
		return filepath.Join(home, ".config", "fish", "config.fish"), nil

	case PowerShell:
		return powerShellProfile(home), nil

	default:
		return "", fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellNoConfigFile, sh), style.ColorError, false))
	}
}

// ----------------------------------
//
//	Resolves the correct $PROFILE path for the running PowerShell flavour.
//	Trusts $PROFILE if it is already set by the active PS session; otherwise
//	derives the path from the OS and available binary:
//	  - pwsh 7+ on Windows        → ~/Documents/PowerShell/Microsoft.PowerShell_profile.ps1
//	  - Windows PowerShell 5.x    → ~/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1
//	  - pwsh on Linux / macOS     → ~/.config/powershell/Microsoft.PowerShell_profile.ps1
//
// ----------------------------------
func powerShellProfile(home string) string {
	// If $PROFILE is already set by the running PS session, trust it.
	if p := os.Getenv("PROFILE"); p != "" {
		return p
	}
	switch runtime.GOOS {
	case "windows":
		// Prefer pwsh 7+ profile; fall back to Windows PowerShell 5.x.
		if _, err := exec.LookPath("pwsh.exe"); err == nil {
			return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
		}
		return filepath.Join(home, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
	default:
		return filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1")
	}
}

// ----------------------------------
//
//	Returns the shell wrapper function snippet to append to the shell config file.
//	The function is always named "rift" and shadows the binary of the same name.
//	To avoid infinite recursion each shell uses its own mechanism to invoke the
//	real binary directly rather than the function:
//	  - bash / zsh / ksh : `command rift`                             (POSIX — skips functions and aliases)
//	  - fish             : `command rift`                             (fish built-in equivalent)
//	  - PowerShell Win   : `rift.exe`                                 (extension disambiguates from the function)
//	  - PowerShell Unix  : `Get-Command -CommandType Application rift`
//
// ----------------------------------
func FunctionSnippet(sh Shell) string {
	switch sh {
	case Fish:
		return `
# rift shell integration
function rift
    eval (command rift $argv)
end
`
	case Ksh:
		return `
# rift shell integration
rift() { eval "$(command rift "$@")"; }
`
	case PowerShell:
		if runtime.GOOS == "windows" {
			return `
# rift shell integration
function rift { Invoke-Expression (rift.exe $args) }
`
		}
		return `
# rift shell integration
function rift { Invoke-Expression (& (Get-Command -CommandType Application rift) $args) }
`
	default:
		// bash and zsh share identical POSIX function syntax.
		return `
# rift shell integration
rift() { eval "$(command rift "$@")"; }
`
	}
}

// ----------------------------------
//
//	Returns the shell command the user should run to reload their config file
//	after shell integration is installed. PowerShell uses `. $PROFILE`;
//	all other shells use `source <cfgFile>`.
//
// ----------------------------------
func ReloadHint(sh Shell, cfgFile string) string {
	if sh == PowerShell {
		return ". $PROFILE"
	}
	return "source " + cfgFile
}

const marker = "# rift shell integration"

// ----------------------------------
//
//	Reports whether the rift shell integration marker is already present in the
//	given config file. Returns false (not an error) if the file does not exist yet.
//
// ----------------------------------
func IsInstalled(configFile string) (bool, error) {
	data, err := os.ReadFile(configFile)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return strings.Contains(string(data), marker), nil
}

// ----------------------------------
//
//	Appends the rift shell function snippet to the shell's config file.
//	Returns an error immediately for CMD and Unsupported shells.
//	Creates the config file's parent directory if it does not exist (relevant
//	for fish and PowerShell on Linux/macOS). Idempotent — skips the write and
//	logs a message if the marker is already present.
//
// ----------------------------------
func Install(sh Shell) error {
	switch sh {
	case CMD:
		return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.ShellCMDNotSupported, style.ColorError, false))
	case Unsupported:
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellUnsupported, os.Getenv("SHELL")), style.ColorError, false))
	}

	cfgFile, err := ConfigFile(sh)
	if err != nil {
		return err
	}

	// Ensure parent directory exists (relevant for fish and PowerShell on Linux/macOS).
	if err := os.MkdirAll(filepath.Dir(cfgFile), 0755); err != nil {
		return err
	}

	installed, err := IsInstalled(cfgFile)
	if err != nil {
		return err
	}
	if installed {
		message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellAlreadyInstalled, cfgFile), style.ColorGreenSoft, false)
		logger.LOGGER.LogToTerminal([]string{message})
		return nil
	}

	f, err := os.OpenFile(cfgFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(FunctionSnippet(sh)); err != nil {
		return err
	}

	message := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellInstallSuccess, cfgFile), style.ColorGreenSoft, false)
	logger.LOGGER.LogToTerminal([]string{message})
	message = style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ShellReloadHint, ReloadHint(sh, cfgFile)), style.ColorCyanSoft, false)
	logger.LOGGER.LogToTerminal([]string{message})
	return nil
}

// ----------------------------------
//
//	Reports whether the rift binary is accessible via PATH.
//
// ----------------------------------
func BinaryInPath() bool {
	_, err := exec.LookPath("rift")
	return err == nil
}
