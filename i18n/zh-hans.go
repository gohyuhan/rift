package i18n

var zH_HANS = LanguageMapping{
	// General
	ConfigPathError:          "配置路径无效，[ERROR: %s]",
	RiftReservedKeywordError: "`%s` 是 rift 的保留关键字",
	RiftDetectedShell:        "rift：检测到的 Shell：%s",
	CWDIsNotDirError:         "当前工作目录不是一个有效的目录",
	PathNotAbsoluteError:     "路径必须为绝对路径，收到: %s",
	NotFileOrDirError:        "指定路径不存在（非文件或目录）",

	// Settings related
	SettingsPathError:                 "无法访问设置目录，[ERROR: %s]",
	SettingsReadError:                 "读取设置文件失败，[ERROR: %s]",
	SettingsParseError:                "解析设置文件失败，已重置为默认值，[ERROR: %s]",
	SettingsWriteError:                "写入设置文件失败，[ERROR: %s]",
	SettingsLanguageUpdated:           "rift: 语言已设置为 %s",
	SettingsLanguageNotSupported:      "rift: 语言 %q 不受支持（支持的语言: %s）",
	SettingsAutoUpdateUpdated:         "rift: 自动更新已设置为 %t",
	SettingsDownloadPreReleaseUpdated: "rift: 下载预发布版本已设置为 %t",

	// DB related
	DBPathError:                 "无法访问数据库目录，[ERROR: %s]",
	DBSetupError:                "数据库初始化失败，[ERROR: %s]",
	DBOpenError:                 "数据库打开失败，请确认是否已初始化，运行 `rift awaken` 进行初始化",
	SettingsBucketNotFoundError: "在数据库中找不到设置存储区，请重新运行 `rift awaken`",
	WaypointBucketNotFoundError: "在数据库中找不到航点存储区，请重新运行 `rift awaken`",
	WaypointDataCorruptedError:  "航点 %q 的数据已损坏，无法读取",

	// Updater related
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

	// Shell related
	ShellCMDNotSupported:  "Windows 命令提示符不支持 Shell 函数。\n请改用 PowerShell、Git Bash 或 WSL，然后重新运行 `rift awaken`。",
	ShellUnsupported:      "rift 不支持 Shell %q。\n支持的 Shell：bash、zsh、fish、ksh、PowerShell。\n您可以手动添加集成 — 请参阅 docs/shell-integration.md",
	ShellNoConfigFile:     "Shell %q 没有已知的配置文件",
	ShellAlreadyInstalled: "rift：Shell 集成已存在于 %s",
	ShellInstallSuccess:   "rift：Shell 集成已添加至 %s",
	ShellInstallReload:    "rift：请重新启动 Shell 或运行：%s",
	BinaryNotInPath:       "rift：在 PATH 中找不到可执行文件 — 请将 rift 添加到 PATH 以便在此会话后继续使用",

	// cmd description
	RiftDescription:                       "通过您预先定义的检查点名称轻松导航路径",
	RiftAwakenDescription:                 "在您的 Shell 中唤醒 rift【首次使用前请执行此命令进行设置与初始化】",
	RiftDiscoverDescription:               "为当前工作目录指定一个航点名称",
	RiftWaypointDescription:               "列出所有航点或显示指定航点的详细信息",
	RiftFlagLanguageDescription:           "设置 rift 的语言（支持：EN、JA、ZH-HANS、ZH-HANT）",
	RiftFlagAutoUpdateDescription:         "设置 rift 是否自动检查更新（使用 --autoupdate 设为启用，--autoupdate=false 设为禁用）",
	RiftFlagDownloadPreReleaseDescription: "设置 rift 是否也下载预发布版本，或仅限稳定版本（使用 --download-pre-release 设为启用，--download-pre-release=false 设为禁用）",
	RiftFlagWaypointDestroyDescription:    "按名称删除一个航点",
	RiftFlagWaypointRebindDescription:     "将现有航点重新绑定到新路径；未提供路径时默认使用当前工作目录，提供有效绝对路径时优先使用该路径",
	RiftFlagWaypointReforgeDescription:    "将现有航点重命名为新名称",

	// Waypoint related
	RiftSavedWaypoint:                     "rift：已保存 %q -> %s",
	RiftUnknownWaypoint:                   "rift：未知的航点名称 %q",
	RiftWaypointAlreadyExistsError:        "航点 %q 已存在，指向 %s",
	RiftWaypointDoNotExistsError:          "rift：航点 %q 不存在",
	RiftWaypointUpdateError:               "rift：更新航点 %q 失败",
	RiftWaypointSealedError:               "rift：航点 %q 已封印，无法前往，原因：%q",
	RiftWaypointSealedLabel:               "（已封印）",
	RiftWaypointRetrieveAllError:          "rift：获取航点列表失败",
	RiftWaypointDetailName:                "航点名称：",
	RiftWaypointDetailPath:                "航点路径：",
	RiftWaypointDetailDiscovered:          "航点发现时间：",
	RiftWaypointDetailTravelledCount:      "航点移动次数：",
	RiftWaypointDetailSealed:              "航点封印：",
	RiftWaypointDetailSealedReason:        "封印原因：",
	RiftWaypointDetailSealedTrue:          "是",
	RiftWaypointDetailSealedFalse:         "否",
	RiftWaypointDestroySuccess:            "rift：航点 %q 已被销毁",
	RiftWaypointDestroyError:              "rift：销毁航点 %q 失败，[ERROR: %s]",
	RiftWaypointRebindNotDirError:         "rift：重新绑定的路径 %q 不是一个目录",
	RiftWaypointRebindSuccess:             "rift：航点 %q 已重新绑定至 %s",
	RiftWaypointRebindError:               "rift：重新绑定航点 %q 失败，[ERROR: %s]",
	RiftWaypointReforgeEmptyError:         "rift：重铸名称不能为空",
	RiftWaypointReforgeError:              "rift：重铸航点 %q 失败，[ERROR: %s]",
	RiftWaypointReforgeAlreadyExistsError: "rift：航点 %q 已存在，无法重铸为已有名称",
	RiftWaypointReforgeSuccess:            "rift：航点 %q 已重铸为 %q",

	// Flag related
	RiftFlagRetrieveError: "rift：获取标志 %q 失败，[ERROR: %s]",

	// Setup related
	CheckAndRunSetupError:  "rift：设置失败，[ERROR: %s]",
	RiftAutoSetupTriggered: "rift：设置与配置已自动触发",
}
