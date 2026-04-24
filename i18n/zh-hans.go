package i18n

var zH_HANS = LanguageMapping{
	// General
	ConfigPathError:          "配置路径无效，[ERROR: %s]",
	RiftReservedKeywordError: "`%s` 是 rift 的保留关键字",
	RiftDetectedShell:        "rift：检测到的 Shell：%s",
	CWDIsNotDirError:         "当前工作目录不是一个有效的目录",
	PathNotAbsoluteError:     "路径必须为绝对路径，收到: %s",
	NotFileOrDirError:        "指定路径不存在（非文件或目录）",
	InvalidValueProvided:     "提供的值无效，不允许包含空格 且不能为空",
	SkippingDueToExecutorErr: "rift：执行器启动失败，符文命令已跳过",

	// Settings
	SettingsPathError:                 "无法访问设置目录，[ERROR: %s]",
	SettingsReadError:                 "读取设置文件失败，[ERROR: %s]",
	SettingsParseError:                "解析设置文件失败，已重置为默认值，[ERROR: %s]",
	SettingsWriteError:                "写入设置文件失败，[ERROR: %s]",
	SettingsLanguageUpdated:           "rift: 语言已设置为 %s",
	SettingsLanguageNotSupported:      "rift: 语言 %q 不受支持（支持的语言: %s）",
	SettingsAutoUpdateUpdated:         "rift: 自动更新已设置为 %t",
	SettingsDownloadPreReleaseUpdated: "rift: 下载预发布版本已设置为 %t",

	// Database
	DBPathError:                 "无法访问数据库目录，[ERROR: %s]",
	DBSetupError:                "数据库初始化失败，[ERROR: %s]",
	DBOpenError:                 "数据库打开失败 — 数据库可能被上一个未正常退出的进程锁定，或有另一个 rift 进程正在运行。若 rift 曾崩溃或被强制终止，请运行 `lsof | grep rift.db` 找到并终止残留进程后重试。若为全新安装，请运行 `rift awaken` 进行初始化",
	WaypointBucketNotFoundError: "在数据库中找不到航点存储区，请重新运行 `rift awaken`",
	WaypointDataCorruptedError:  "航点 [%s] 的数据已损坏，无法读取",
	SpellBucketNotFoundError:    "在数据库中找不到咒语存储区，请重新运行 `rift awaken`",
	SpellDataCorruptedError:     "咒语 [%s] 的数据已损坏，无法读取",
	RuneBucketNotFoundError:     "在数据库中找不到符文存储区，请重新运行 `rift awaken`",
	RuneDataCorruptedError:      "路径 [%s] 的符文数据已损坏，无法读取",

	// Updater
	UpdaterDownloadPrompt:               "发现新版本 %s，是否立即下载？(y/n): ",
	UpdaterFailToCheckForUpdate:         "检查更新失败：%v",
	UpdaterFailToCreateRequest:          "无法创建请求：%v",
	UpdaterFailToFetchRelease:           "无法获取最新版本信息：%v",
	UpdaterNoReleasesFound:              "未找到任何版本",
	UpdaterFailToReadResponse:           "无法读取响应内容：%v",
	UpdaterFailToParseJSON:              "无法解析 JSON 响应：%v",
	UpdaterFailToExtractBinary:          "无法解压可执行文件：%v",
	UpdaterUnsupportedArchiveFormat:     "不支持的压缩格式",
	UpdaterBinaryNotFoundInArchive:      "在压缩包中找不到可执行文件",
	UpdaterAlreadyLatest:                "您已是最新版本（%s）",
	UpdaterDownloading:                  "正在下载版本 %s...",
	UpdaterUnSupportedOS:                "不支持的操作系统/架构：%s/%s",
	UpdaterDownloadFail:                 "下载更新失败：%v",
	UpdaterBinaryReplaceFail:            "替换可执行文件失败：%v",
	UpdaterDownloadSuccess:              "成功更新至版本 %s",
	UpdaterDownloadUnexpectedStatusCode: "意外的状态码：%d",
	UpdaterRequiresSudo:                 "权限不足，尝试以 sudo 重试...",

	// Shell
	ShellCMDNotSupported:  "Windows 命令提示符不支持 Shell 函数。\n请改用 PowerShell、Git Bash 或 WSL，然后重新运行 `rift awaken`。",
	ShellUnsupported:      "rift 不支持 Shell %q。\n支持的 Shell：bash、zsh、fish、ksh、PowerShell。\n您可以手动添加集成 — 请参阅 docs/shell-integration.md",
	ShellNoConfigFile:     "Shell %q 没有已知的配置文件",
	ShellAlreadyInstalled: "rift：Shell 集成已存在于 %s",
	ShellInstallSuccess:   "rift：Shell 集成已添加至 %s",
	ShellReloadHint:       "rift：请重新启动 Shell 或运行：%s",
	BinaryNotInPath:       "rift：在 PATH 中找不到可执行文件 — 请将 rift 添加到 PATH 以便在此会话后继续使用",

	// Commands and flags
	RiftDescription:                       "通过您预先定义的检查点名称轻松导航路径",
	RiftAwakenDescription:                 "在您的 Shell 中唤醒 rift【首次使用前请执行此命令进行设置与初始化】",
	RiftDiscoverDescription:               "为当前工作目录指定一个航点名称",
	RiftWaypointDescription:               "启动航点交互界面或显示指定航点的详细信息",
	RiftLearnDescription:                  "通过为指令指定名称来教 rift 一个新咒语；多词指令请用引号括起（例：rift learn build \"docker compose up --build\"）",
	RiftSpellDescription:                  "通过已习得的咒语名称施咒，或直接执行命令字符串（例：\"git commit -m 'msg'\"）",
	RiftSpellbookDescription:              "启动咒语交互界面或显示指定咒语的详细信息",
	RiftFlagLanguageDescription:           "设置 rift 的语言（支持：EN、JA、ZH-HANS、ZH-HANT）",
	RiftFlagAutoUpdateDescription:         "设置 rift 是否自动检查更新（使用 --autoupdate 设为启用，--autoupdate=false 设为禁用）",
	RiftFlagDownloadPreReleaseDescription: "设置 rift 是否也下载预发布版本，或仅限稳定版本（使用 --download-pre-release 设为启用，--download-pre-release=false 设为禁用）",
	RiftFlagWaypointDestroyDescription:    "按名称删除一个航点",
	RiftFlagSpellForgetDescription:        "按名称删除一个已保存的咒语",
	RiftFlagWaypointRebindDescription:     "将现有航点重新绑定到新路径；未提供路径时默认使用当前工作目录，提供有效绝对路径时优先使用该路径",
	RiftFlagWaypointReforgeDescription:    "将现有航点重命名为新名称",
	RiftFlagUpdateDescription:             "手动触发检查最新版本，如有可用更新则进行升级",
	RiftFlagVersionDescription:            "打印 rift 的当前版本",
	RiftFlagCastDescription:               "代替导航，在航点路径下施放已习得的咒语或执行命令字符串（例：\"git commit -m 'msg'\"）；命令将以航点路径为工作目录执行",
	RiftFlagRetrieveError:                 "rift：获取标志 %q 失败，[ERROR: %s]",
	RiftRuneDescription:                   "为航点绑定进入和离开时的触发命令；当 rift 导航至或离开该航点时自动执行",

	// Spell operations
	RiftSpellSaved:            "rift：已习得 %q -> %s",
	RiftSpellUpdated:          "rift：咒语 %q 已更新 -> %s",
	RiftSpellForgetSuccess:    "rift：咒语 %q 已被遗忘",
	RiftSpellForgetError:      "rift：遗忘咒语 %q 失败，[ERROR: %s]",
	RiftSpellDoNotExistsError: "rift：咒语 %q 不存在",
	RiftSpellUpdateError:      "rift：更新咒语 %q 失败，[ERROR: %s]",
	ForbiddenShellBuiltinSpellCommand: "rift：Shell 内建命令（如 cd、export、source、alias）只影响执行它们的进程，无法修改当前 Shell 会话。" +
		"如需在命令序列中使用内建命令，请使用 Shell 的 -c 参数显式调用 Shell 并链接命令。" +
		"使用 --login（或等效选项）可在该进程中加载完整的 Shell 环境（PATH、别名、配置文件等）。\n\n" +
		"示例:\n" +
		"  bash --login -c \"source env/bin/activate && python main.py\"\n" +
		"  zsh --login -c \"source env/bin/activate && python main.py\"\n" +
		"  fish --login -c \"source env/bin/activate.fish; python main.py\"\n" +
		"  pwsh -Login -Command \". ./env/bin/Activate.ps1; python main.py\"\n" +
		"  cmd /c \"activate.bat && python main.py\"",
	ForbiddenRiftNavigationSpellCommand: "rift：rift 路径导航命令（如 rift <路径点名称>）不能作为咒语学习 — 它们由 Shell 集成直接处理，不会作为子进程运行。",
	ForbiddenRiftNavigationRuneCommand:  "rift：rift 路径导航命令（如 rift <路径点名称>）不支持在符文中使用 — 如需在特定目录运行命令，请改用 `rift <路径点名称> --cast <咒语名称>`。",
	SpellCommandEmpty:                   "rift：咒语命令不能为空",
	InvalidSpellCommandError:            "rift：咒语命令 [%s] 无效，无法执行",

	// Waypoint operations
	RiftSavedWaypoint:                     "rift：已保存 %q -> %s",
	RiftWaypointAlreadyExistsError:        "航点 %q 已存在，指向 %s",
	RiftWaypointDoNotExistsError:          "rift：航点 %q 不存在",
	RiftWaypointUpdateError:               "rift：更新航点 %q 失败，[ERROR: %s]",
	RiftWaypointSealedError:               "rift：航点 %q 已封印，无法前往，原因：%q",
	RiftWaypointSealedLabel:               "（已封印）",
	RiftWaypointRetrieveAllError:          "rift：获取航点列表失败：%s",
	RiftWaypointDestroySuccess:            "rift：航点 %q 已被销毁",
	RiftWaypointDestroyError:              "rift：销毁航点 %q 失败，[ERROR: %s]",
	RiftWaypointRebindNotDirError:         "rift：重新绑定的路径 %q 不是一个目录",
	RiftWaypointRebindSuccess:             "rift：航点 %q 已重新绑定至 %s",
	RiftWaypointReforgeEmptyError:         "rift：重铸名称不能为空",
	RiftWaypointReforgeAlreadyExistsError: "rift：航点 %q 已存在，无法重铸为已有名称",
	RiftWaypointReforgeError:              "rift：重铸航点 %q 失败，[ERROR: %s]",
	RiftWaypointReforgeSuccess:            "rift：航点 %q 已重铸为 %q",

	// Rune operations
	RiftRuneEngraveSuccessful: "rift：符文已刻印到航点 %q",
	RiftRuneEngraveNone:       "rift：未将符文刻印到航点 %q",
	RiftRuneUpdateError:       "rift：更新路径 %q 的符文失败，[ERROR: %s]",

	// Spell detail view
	RiftSpellDetailName:      "咒语名称：",
	RiftSpellDetailCommand:   "咒语命令：",
	RiftSpellDetailAddedAt:   "咒语添加时间：",
	RiftSpellDetailCastCount: "咒语施放次数：",

	// Waypoint detail view
	RiftWaypointDetailName:           "航点名称：",
	RiftWaypointDetailPath:           "航点路径：",
	RiftWaypointDetailDiscovered:     "航点发现时间：",
	RiftWaypointDetailTravelledCount: "航点移动次数：",
	RiftWaypointDetailSealed:         "航点封印：",
	RiftWaypointDetailSealedReason:   "封印原因：",
	RiftWaypointDetailSealedTrue:     "是",
	RiftWaypointDetailSealedFalse:    "否",

	// Waypoint interactive UI
	WaypointInfoListTitle:                         "航点列表",
	WaypointInteractiveError:                      "rift：航点交互界面启动失败，[ERROR: %s]",
	RebindPathInputPlaceHolder:                    "输入绝对路径，留空则使用当前工作目录",
	WaypointRebindTitle:                           "重新绑定",
	ReforgeWaypointNameInputPlaceHolder:           "输入航点的新名称",
	WaypointReforgeTitle:                          "重铸",
	WaypointUIUpKeyHelp:                           "上移",
	WaypointUIUpKeyHelpDescription:                "将光标移至上一个航点",
	WaypointUIDownKeyHelp:                         "下移",
	WaypointUIDownKeyHelpDescription:              "将光标移至下一个航点",
	WaypointUIQuitKeyHelp:                         "退出",
	WaypointUIQuitKeyHelpDescription:              "退出航点交互界面",
	WaypointUIHelpKeyHelp:                         "帮助",
	WaypointUIHelpKeyHelpDescription:              "显示或隐藏完整的快捷键列表",
	WaypointNavigateKeyHelp:                       "导航",
	WaypointNavigateKeyHelpDescription:            "前往所选航点",
	WaypointDestroyKeyHelp:                        "销毁",
	WaypointDestroyKeyHelpDescription:             "永久删除所选航点",
	WaypointUnsealKeyHelp:                         "解封",
	WaypointUnsealKeyHelpDescription:              "尝试解封所选航点",
	WaypointRebindKeyHelp:                         "重新绑定航点路径",
	WaypointRebindKeyHelpDescription:              "将所选航点重新分配至新路径",
	WaypointReforgeKeyHelp:                        "重铸航点名称",
	WaypointReforgeKeyHelpDescription:             "将所选航点重命名为新名称",
	WaypointNameCopyPathCopyKeyHelp:               "复制航点名称 / 路径到剪贴板",
	WaypointNameCopyPathCopyKeyHelpDescription:    "将航点名称 (y) 或路径 (Y) 复制到剪贴板",
	WaypointCopyFromInputValueKeyHelp:             "复制输入内容",
	WaypointCopyFromInputValueKeyHelpDescription:  "将当前输入框的内容复制到剪贴板",
	WaypointPasteIntoInputValueKeyHelp:            "粘贴至输入框",
	WaypointPasteIntoInputValueKeyHelpDescription: "将剪贴板内容粘贴至输入框",
	WaypointClosePopUp:                            "关闭弹窗",
	WaypointClosePopUpDescription:                 "关闭当前弹窗而不保存",

	// Spell interactive UI
	SpellInfoListTitle:                       "法术书",
	SpellbookInteractiveError:                "rift：法术书交互界面启动失败，[ERROR: %s]",
	RiftSpellRetrieveAllError:                "rift：获取法术失败",
	SpellUIUpKeyHelp:                         "上移",
	SpellUIUpKeyHelpDescription:              "将光标移至上一个法术",
	SpellUIDownKeyHelp:                       "下移",
	SpellUIDownKeyHelpDescription:            "将光标移至下一个法术",
	SpellUIQuitKeyHelp:                       "退出",
	SpellUIQuitKeyHelpDescription:            "退出法术书交互界面",
	SpellUIHelpKeyHelp:                       "帮助",
	SpellUIHelpKeyHelpDescription:            "显示或隐藏完整的快捷键列表",
	SpellCastKeyHelp:                         "施放",
	SpellCastKeyHelpDescription:              "在当前工作目录施放法术，或施放至已发现的航点",
	SpellLearnKeyHelp:                        "习得",
	SpellLearnKeyHelpDescription:             "习得一个新法术",
	SpellForgetKey:                           "遗忘",
	SpellForgetKeyDescription:                "永久遗忘所选法术",
	SpellClosePopUp:                          "关闭弹窗",
	SpellClosePopUpDescription:               "关闭当前弹窗而不保存",
	SpellUILearnKeyHelp:                      "确认",
	SpellUINextInputKeyHelp:                  "下一个输入",
	SpellUIPreviousInputKeyHelp:              "上一个输入",
	SpellNameInputPlaceHolder:                "输入咒语的名称",
	SpellCommandInputPlaceHolder:             "输入咒语的命令",
	SpellNameInputTitle:                      "咒语名称：",
	SpellCommandInputTitle:                   "咒语命令：",
	SpellUIChooseCastLocationKeyHelp:         "选择施放位置",
	SpellUIChooseWaypointCastLocationKeyHelp: "选择施放咒语的航点",

	// Rune interactive UI
	RuneInteractiveError:           "[ERROR: %s]",
	RuneEngraveTypeOptionListTitle: "符文选项",
	EngraveRuneEnterTitle:          "进入时符文的命令：",
	EngraveRuneLeaveTitle:          "离开时符文的命令：",
	EngraveRuneEngraveButton:       "刻印",
	RuneCommandsPlaceHolder:        "输入命令…（cd 无效，推荐使用 rift 切换路径）",
	RuneCommandsInvalidDueToShellBuildInCommand: "检测到符文使用了 Shell 内建命令（如 cd、export、source、alias）——内建命令只影响执行它们的进程，无法修改当前 Shell 会话。" +
		"如需在命令序列中使用内建命令，请使用 Shell 的 -c 参数显式调用 Shell 并链接命令。" +
		"使用 -i（交互模式）可加载 Shell 的交互配置（.zshrc、.bashrc 等）——若环境依赖其中内容（如使用 rift 切换目录、nvm、conda）则必须使用此选项。\n\n" +
		"示例:\n" +
		"  bash -i -c \"source env/bin/activate && python main.py\"\n" +
		"  zsh -i -c \"source env/bin/activate && python main.py\"\n" +
		"  fish -i -c \"source env/bin/activate.fish; python main.py\"\n" +
		"  pwsh -Login -Command \". ./env/bin/Activate.ps1; python main.py\"\n" +
		"  cmd /c \"activate.bat && python main.py\"",
	EngraveRuneEnterOptionName: "刻印进入时符文",
	EngraveRuneEnterOptionDesc: "设置进入此航点时执行的命令",
	EngraveRuneLeaveOptionName: "刻印离开时符文",
	EngraveRuneLeaveOptionDesc: "设置离开此航点时执行的命令",
	RemoveRuneEnterOptionName:  "移除进入时符文",
	RemoveRuneEnterOptionDesc:  "清除进入此航点时执行的命令",
	RemoveRuneLeaveOptionName:  "移除离开时符文",
	RemoveRuneLeaveOptionDesc:  "清除离开此航点时执行的命令",

	// Cast location option popup
	CastLocationOptionTitle:               "施放位置",
	CastLocationOptionCurrent:             "当前目录",
	CastLocationOptionCurrentDescription:  "在当前工作目录中施放咒语",
	CastLocationOptionWaypoint:            "航点",
	CastLocationOptionWaypointDescription: "从已发现的航点中选择并在那里施放咒语",
	CastWaypointLocationOptionTitle:       "选择航点",

	// Setup
	CheckAndRunSetupError:  "rift：设置失败，[ERROR: %s]",
	RiftAutoSetupTriggered: "rift：设置与配置已自动触发",
}
