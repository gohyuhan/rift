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
var APPVERSION = "v0.2.0-pr.2"

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
