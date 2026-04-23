package constant

const (
	APPNAME         = "rift"
	APPDBNAME       = "rift.db"
	APPSETTINGSNAME = "rift_settings.json"

	APPLOGO = `
  ██████╗ ██╗███████╗████████╗
  ██╔══██╗██║██╔════╝╚══██╔══╝
  ██████╔╝██║█████╗     ██║
  ██╔══██╗██║██╔══╝     ██║
  ██║  ██║██║██║        ██║
  ╚═╝  ╚═╝╚═╝╚═╝        ╚═╝
 ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
 ┃  Tear through your filesystem.  ┃
 ┃  No spells required.            ┃
 ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛`
)

// this will be injected during build
// example: go build -ldflags "-X rift/constant.APPVERSION=v1.x.x" -o main
var APPVERSION = "v0.3.0"

const (
	RIFT_CMD_KEYWORD      = "rift"
	AWAKEN_CMD_KEYWORD    = "awaken"
	DISCOVER_CMD_KEYWORD  = "discover"
	WAYPOINT_CMD_KEYWORD  = "waypoint"
	SPELL_CMD_KEYWORD     = "spell"
	LEARN_CMD_KEYWORD     = "learn"
	SPELLBOOK_CMD_KEYWORD = "spellbook"
	CAST_CMD_KEYWORD      = "cast"
	RITUAL_CMD_KEYWORD    = "ritual"
	INSCRIBE_CMD_KEYWORD  = "inscribe"
	SCROLL_CMD_KEYWORD    = "scroll"
	SORCERY_CMD_KEYWORD   = "sorcery"
	SUMMON_CMD_KEYWORD    = "summon"
	DEPLOY_CMD_KEYWORD    = "deploy"
	RUNE_CMD_KEYWORD      = "rune"
	SEER_CMD_KEYWORD      = "seer"
	RECALL_CMD_KEYWORD    = "recall"
	LOOT_CMD_KEYWORD      = "loot"
	GRIMOIRE_CMD_KEYWORD  = "grimoire"
	LORE_CMD_KEYWORD      = "lore"
	STATS_CMD_KEYWORD     = "stats"
)

// First-argument commands that mutate parent-shell state (cwd, env, aliases,
// functions, options, job table, modules) and therefore cannot be run via
// exec.Command — a subprocess change dies with the subprocess. Route these
// through executor.RunShellBuiltInCmd so the shell wrapper evals them in the
// caller's shell context. Deduplicated: each token appears exactly once in
// the section where it is most characteristic.
var ShellBuildInCmd = []string{
	// ------------------------------------------------------------------
	// General — POSIX, works in bash / zsh / ksh / fish / PowerShell
	// ------------------------------------------------------------------
	"cd",
	"pwd",
	"export",
	"unset",
	"alias",
	"unalias",
	"source",
	".",
	"eval",
	"exec",
	"set",
	"shift",
	"umask",
	"ulimit",
	"trap",
	"exit",
	"logout",
	"read",
	"getopts",
	"suspend",

	// ------------------------------------------------------------------
	// bash / zsh / ksh shared
	// ------------------------------------------------------------------
	"local",
	"declare",
	"typeset",
	"readonly",
	"nameref",
	"pushd",
	"popd",
	"dirs",
	"hash",
	"history",
	"fc",
	"jobs",
	"bg",
	"fg",
	"disown",
	"wait",
	"enable",
	"disable",

	// ------------------------------------------------------------------
	// bash specific
	// ------------------------------------------------------------------
	"bind",
	"complete",
	"compgen",
	"compopt",
	"shopt",
	"caller",
	"mapfile",
	"readarray",
	"let",

	// ------------------------------------------------------------------
	// zsh specific
	// ------------------------------------------------------------------
	"setopt",
	"unsetopt",
	"autoload",
	"zmodload",
	"bindkey",
	"compdef",
	"zstyle",
	"emulate",
	"vared",
	"limit",
	"unlimit",
	"sched",

	// ------------------------------------------------------------------
	// ksh specific
	// ------------------------------------------------------------------
	"whence",
	"print",

	// ------------------------------------------------------------------
	// fish specific
	// ------------------------------------------------------------------
	"function",
	"functions",
	"funced",
	"funcsave",
	"abbr",
	"fish_add_path",
	"status",
	"emit",
	"block",

	// ------------------------------------------------------------------
	// PowerShell specific
	// ------------------------------------------------------------------
	"Set-Location",
	"Push-Location",
	"Pop-Location",
	"Set-Variable",
	"Remove-Variable",
	"Clear-Variable",
	"Set-Alias",
	"New-Alias",
	"Remove-Alias",
	"Clear-Alias",
	"Import-Module",
	"Remove-Module",
	"New-PSDrive",
	"Remove-PSDrive",
	"Add-Type",
	"Enter-PSSession",
	"Exit-PSSession",
	"Update-FormatData",
	"Update-TypeData",
	"Register-ArgumentCompleter",

	// ------------------------------------------------------------------
	// Sourced activation scripts — shell functions injected into parent
	// (venvs, version managers). Not builtins; lost across subprocess.
	// ------------------------------------------------------------------
	"activate",
	"deactivate",
	"workon",
	"nvm",
	"pyenv",
	"rbenv",
	"jenv",
	"goenv",
	"gvm",
	"sdk",
	"asdf",
	"fnm",
	"volta",
	"rustup",
	"conda",
	"mamba",
	"direnv",

	// ------------------------------------------------------------------
	// Third-party shell function shims — installed by user's rc file
	// ------------------------------------------------------------------
	"z",
	"zi",
	"zshz",
	"j",
	"jc",
	"fasd",
	"kubectx",
	"kubens",
	"starship",
	"thefuck",
	"fuck",
}
