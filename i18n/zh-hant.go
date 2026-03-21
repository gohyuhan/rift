package i18n

var zH_HANT = LanguageMapping{
	// General
	ConfigPathError:          "設定路徑無效，[ERROR: %s]",
	RiftReservedKeywordError: "`%s` 是 rift 的保留關鍵字",
	RiftDetectedShell:        "rift：偵測到的 Shell：%s",
	CWDIsNotDirError:         "目前工作目錄不是一個有效的目錄",
	PathNotAbsoluteError:     "路徑必須為絕對路徑，收到: %s",
	NotFileOrDirError:        "指定路徑不存在（非檔案或目錄）",

	// Settings related
	SettingsPathError:                 "無法存取設定目錄，[ERROR: %s]",
	SettingsReadError:                 "讀取設定檔失敗，[ERROR: %s]",
	SettingsParseError:                "解析設定檔失敗，已重置為預設值，[ERROR: %s]",
	SettingsWriteError:                "寫入設定檔失敗，[ERROR: %s]",
	SettingsLanguageUpdated:           "rift: 語言已設定為 %s",
	SettingsLanguageNotSupported:      "rift: 語言 %q 不受支援（支援的語言: %s）",
	SettingsAutoUpdateUpdated:         "rift: 自動更新已設定為 %t",
	SettingsDownloadPreReleaseUpdated: "rift: 下載預發布版本已設定為 %t",

	// DB related
	DBPathError:                 "無法存取資料庫目錄，[ERROR: %s]",
	DBSetupError:                "資料庫初始化失敗，[ERROR: %s]",
	DBOpenError:                 "資料庫開啟失敗，請確認是否已初始化，執行 `rift awaken` 進行初始化",
	SettingsBucketNotFoundError: "在資料庫中找不到設定儲存區，請重新執行 `rift awaken`",
	WaypointBucketNotFoundError: "在資料庫中找不到航點儲存區，請重新執行 `rift awaken`",
	WaypointDataCorruptedError:  "航點 %q 的資料已損毀，無法讀取",

	// Updater related
	UpdaterDownloadPrompt:               "發現新版本 %s，是否立即下載？(y/n): ",
	UpdaterFailToCheckForUpdate:         "檢查更新失敗：%v",
	UpdaterFailToCreateRequest:          "無法建立請求：%v",
	UpdaterFailToFetchRelease:           "無法取得最新版本資訊：%v",
	UpdaterNoReleasesFound:              "找不到任何版本",
	UpdaterFailToReadResponse:           "無法讀取回應內容：%v",
	UpdaterFailToParseJSON:              "無法解析 JSON 回應：%v",
	UpdaterFailToExtractBinary:          "無法解壓縮執行檔：%v",
	UpdaterUnsupportedArchiveFormat:     "不支援的壓縮格式",
	UpdaterBinaryNotFoundInArchive:      "在壓縮檔中找不到執行檔",
	UpdaterAlreadyLatest:                "您已是最新版本（%s）",
	UpdaterDownloading:                  "正在下載版本 %s...",
	UpdaterUnSupportedOS:                "不支援的作業系統/架構：%s/%s",
	UpdaterDownloadFail:                 "下載更新失敗：%v",
	UpdaterBinaryReplaceFail:            "替換執行檔失敗：%v",
	UpdaterDownloadSuccess:              "成功更新至版本 %s",
	UpdaterDownloadUnexpectedStatusCode: "非預期的狀態碼：%d",
	UpdaterRequiresSudo:                 "權限不足，嘗試以 sudo 重試...",

	// Shell related
	ShellCMDNotSupported:  "Windows 命令提示字元不支援 Shell 函式。\n請改用 PowerShell、Git Bash 或 WSL，然後重新執行 `rift awaken`。",
	ShellUnsupported:      "rift 不支援 Shell %q。\n支援的 Shell：bash、zsh、fish、ksh、PowerShell。\n您可以手動新增整合 — 請參閱 docs/shell-integration.md",
	ShellNoConfigFile:     "Shell %q 沒有已知的設定檔",
	ShellAlreadyInstalled: "rift：Shell 整合已存在於 %s",
	ShellInstallSuccess:   "rift：Shell 整合已新增至 %s",
	ShellInstallReload:    "rift：請重新啟動 Shell 或執行：%s",
	BinaryNotInPath:       "rift：在 PATH 中找不到執行檔 — 請將 rift 加入 PATH 以便在此工作階段後繼續使用",

	// cmd description
	RiftDescription:                       "透過您預先定義的檢查點名稱輕鬆導航路徑",
	RiftAwakenDescription:                 "在您的 Shell 中喚醒 rift【首次使用前請執行此指令進行設定與初始化】",
	RiftDiscoverDescription:               "為目前工作目錄指定一個航點名稱",
	RiftWaypointDescription:               "列出所有航點或顯示指定航點的詳細資訊",
	RiftFlagLanguageDescription:           "設定 rift 的語言（支援：EN、JA、ZH-HANS、ZH-HANT）",
	RiftFlagAutoUpdateDescription:         "設定 rift 是否自動檢查更新（使用 --autoupdate 設為啟用，--autoupdate=false 設為停用）",
	RiftFlagDownloadPreReleaseDescription: "設定 rift 是否也下載預發布版本，或僅限穩定版本（使用 --download-pre-release 設為啟用，--download-pre-release=false 設為停用）",
	RiftFlagWaypointDestroyDescription:    "依名稱移除一個航點",
	RiftFlagWaypointRebindDescription:     "將現有航點重新綁定至新路徑；未提供路徑時預設使用目前工作目錄，提供有效絕對路徑時優先使用該路徑",
	RiftFlagWaypointReforgeDescription:    "將現有航點重新命名為新名稱",

	// Waypoint related
	RiftSavedWaypoint:                "rift：已儲存 %q -> %s",
	RiftUnknownWaypoint:              "rift：未知的航點名稱 %q",
	RiftWaypointAlreadyExistsError:   "航點 %q 已存在，指向 %s",
	RiftWaypointDoNotExistsError:     "rift：航點 %q 不存在",
	RiftWaypointUpdateError:          "rift：更新航點 %q 失敗",
	RiftWaypointSealedError:          "rift：航點 %q 已封印，無法前往，原因：%q",
	RiftWaypointSealedLabel:          "（已封印）",
	RiftWaypointRetrieveAllError:     "rift：取得航點清單失敗",
	RiftWaypointDetailName:           "航點名稱：",
	RiftWaypointDetailPath:           "航點路徑：",
	RiftWaypointDetailDiscovered:     "航點發現時間：",
	RiftWaypointDetailTravelledCount: "航點移動次數：",
	RiftWaypointDetailSealed:         "航點封印：",
	RiftWaypointDetailSealedReason:   "封印原因：",
	RiftWaypointDetailSealedTrue:     "是",
	RiftWaypointDetailSealedFalse:    "否",

	// Flag related
	RiftFlagRetrieveError: "rift：取得旗標 %q 失敗，[ERROR: %s]",

	// Setup related
	CheckAndRunSetupError:  "rift：設定失敗，[ERROR: %s]",
	RiftAutoSetupTriggered: "rift：設定與配置已自動觸發",
}
