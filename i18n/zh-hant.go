package i18n

var zH_HANT = LanguageMapping{
	// General
	ConfigPathError:          "設定路徑無效，[ERROR: %s]",
	RiftReservedKeywordError: "`%s` 是 rift 的保留關鍵字",
	RiftDetectedShell:        "rift：偵測到的 Shell：%s",
	CWDIsNotDirError:         "目前工作目錄不是一個有效的目錄",
	PathNotAbsoluteError:     "路徑必須為絕對路徑，收到: %s",
	NotFileOrDirError:        "指定路徑不存在（非檔案或目錄）",
	InvalidValueProvided:     "提供的值無效，不允許包含空格 且不能为空",
	SkippingDueToExecutorErr: "rift：執行器啟動失敗，符文命令已略過",

	// Settings
	SettingsPathError:                 "無法存取設定目錄，[ERROR: %s]",
	SettingsReadError:                 "讀取設定檔失敗，[ERROR: %s]",
	SettingsParseError:                "解析設定檔失敗，已重置為預設值，[ERROR: %s]",
	SettingsWriteError:                "寫入設定檔失敗，[ERROR: %s]",
	SettingsLanguageUpdated:           "rift: 語言已設定為 %s",
	SettingsLanguageNotSupported:      "rift: 語言 %q 不受支援（支援的語言: %s）",
	SettingsAutoUpdateUpdated:         "rift: 自動更新已設定為 %t",
	SettingsDownloadPreReleaseUpdated: "rift: 下載預發布版本已設定為 %t",

	// Database
	DBPathError:                 "無法存取資料庫目錄，[ERROR: %s]",
	DBSetupError:                "資料庫初始化失敗，[ERROR: %s]",
	DBOpenError:                 "資料庫開啟失敗 — 資料庫可能被上一個未正常退出的程序鎖定，或另一個 rift 程序正在執行中。若 rift 曾崩潰或被強制結束，請執行 `lsof | grep rift.db` 找到並終止殘留程序後重試。若為全新安裝，請執行 `rift awaken` 進行初始化",
	WaypointBucketNotFoundError: "在資料庫中找不到航點儲存區，請重新執行 `rift awaken`",
	WaypointDataCorruptedError:  "航點 [%s] 的資料已損毀，無法讀取",
	SpellBucketNotFoundError:    "在資料庫中找不到咒語儲存區，請重新執行 `rift awaken`",
	SpellDataCorruptedError:     "咒語 [%s] 的資料已損毀，無法讀取",
	RuneBucketNotFoundError:     "在資料庫中找不到符文儲存區，請重新執行 `rift awaken`",
	RuneDataCorruptedError:      "路徑 [%s] 的符文資料已損毀，無法讀取",

	// Updater
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

	// Shell
	ShellCMDNotSupported:  "Windows 命令提示字元不支援 Shell 函式。\n請改用 PowerShell、Git Bash 或 WSL，然後重新執行 `rift awaken`。",
	ShellUnsupported:      "rift 不支援 Shell %q。\n支援的 Shell：bash、zsh、fish、ksh、PowerShell。\n您可以手動新增整合 — 請參閱 docs/shell-integration.md",
	ShellNoConfigFile:     "Shell %q 沒有已知的設定檔",
	ShellAlreadyInstalled: "rift：Shell 整合已存在於 %s",
	ShellInstallSuccess:   "rift：Shell 整合已新增至 %s",
	ShellReloadHint:       "rift：請重新啟動 Shell 或執行：%s",
	BinaryNotInPath:       "rift：在 PATH 中找不到執行檔 — 請將 rift 加入 PATH 以便在此工作階段後繼續使用",

	// Commands and flags
	RiftDescription:                       "透過您預先定義的檢查點名稱輕鬆導航路徑",
	RiftAwakenDescription:                 "在您的 Shell 中喚醒 rift【首次使用前請執行此指令進行設定與初始化】",
	RiftDiscoverDescription:               "為目前工作目錄指定一個航點名稱",
	RiftWaypointDescription:               "啟動航點互動介面或顯示指定航點的詳細資訊",
	RiftLearnDescription:                  "透過為指令指定名稱來教 rift 一個新咒語；多詞指令請用引號括起（例：rift learn build \"docker compose up --build\"）",
	RiftSpellDescription:                  "透過咒語名稱施咒，執行其綁定的終端命令",
	RiftSpellbookDescription:              "啟動咒語互動介面或顯示指定咒語的詳細資訊",
	RiftFlagLanguageDescription:           "設定 rift 的語言（支援：EN、JA、ZH-HANS、ZH-HANT）",
	RiftFlagAutoUpdateDescription:         "設定 rift 是否自動檢查更新（使用 --autoupdate 設為啟用，--autoupdate=false 設為停用）",
	RiftFlagDownloadPreReleaseDescription: "設定 rift 是否也下載預發布版本，或僅限穩定版本（使用 --download-pre-release 設為啟用，--download-pre-release=false 設為停用）",
	RiftFlagWaypointDestroyDescription:    "依名稱移除一個航點",
	RiftFlagSpellForgetDescription:        "依名稱移除一個已儲存的咒語",
	RiftFlagWaypointRebindDescription:     "將現有航點重新綁定至新路徑；未提供路徑時預設使用目前工作目錄，提供有效絕對路徑時優先使用該路徑",
	RiftFlagWaypointReforgeDescription:    "將現有航點重新命名為新名稱",
	RiftFlagUpdateDescription:             "手動觸發檢查最新版本，如有可用更新則進行升級",
	RiftFlagVersionDescription:            "列印 rift 的目前版本",
	RiftFlagCastDescription:               "取代導航，對航點的路徑施放已習得的咒語；咒語命令將在航點的路徑作為工作目錄時執行",
	RiftFlagRetrieveError:                 "rift：取得旗標 %q 失敗，[ERROR: %s]",
	RiftRuneDescription:                   "為航點綁定進入和離開時的觸發命令；當 rift 導航至或離開該航點時自動執行",

	// Spell operations
	RiftSpellSaved:            "rift：已習得 %q -> %s",
	RiftSpellUpdated:          "rift：咒語 %q 已更新 -> %s",
	RiftSpellForgetSuccess:    "rift：咒語 %q 已被遺忘",
	RiftSpellForgetError:      "rift：遺忘咒語 %q 失敗，[ERROR: %s]",
	RiftSpellDoNotExistsError: "rift：咒語 %q 不存在",
	RiftSpellUpdateError:      "rift：更新咒語 %q 失敗，[ERROR: %s]",
	ForbiddenShellBuiltinSpellCommand: "rift：Shell 內建命令（如 cd、export、source、alias）只影響執行它們的進程，無法修改當前 Shell 會話。" +
		"如需在命令序列中使用內建命令，請使用 Shell 的 -c 參數顯式呼叫 Shell 並串連命令。" +
		"使用 --login（或等效選項）可在該進程中載入完整的 Shell 環境（PATH、別名、設定檔等）。\n\n" +
		"範例:\n" +
		"  bash --login -c \"source env/bin/activate && python main.py\"\n" +
		"  zsh --login -c \"source env/bin/activate && python main.py\"\n" +
		"  fish --login -c \"source env/bin/activate.fish; python main.py\"\n" +
		"  pwsh -Login -Command \". ./env/bin/Activate.ps1; python main.py\"\n" +
		"  cmd /c \"activate.bat && python main.py\"",
	ForbiddenRiftNavigationSpellCommand: "rift：rift 路徑導航命令（如 rift <路徑點名稱>）不能作為咒語學習 — 它們由 Shell 整合直接處理，不會作為子程序執行。",
	ForbiddenRiftNavigationRuneCommand:  "rift：rift 路徑導航命令（如 rift <路徑點名稱>）不支援在符文中使用 — 如需在特定目錄執行命令，請改用 `rift <路徑點名稱> --cast <咒語名稱>`。",
	SpellCommandEmpty:                   "rift：咒語命令不能為空",
	InvalidSpellCommandError:            "rift：咒語命令 [%s] 無效，無法執行",

	// Waypoint operations
	RiftSavedWaypoint:                     "rift：已儲存 %q -> %s",
	RiftWaypointAlreadyExistsError:        "航點 %q 已存在，指向 %s",
	RiftWaypointDoNotExistsError:          "rift：航點 %q 不存在",
	RiftWaypointUpdateError:               "rift：更新航點 %q 失敗，[ERROR: %s]",
	RiftWaypointSealedError:               "rift：航點 %q 已封印，無法前往，原因：%q",
	RiftWaypointSealedLabel:               "（已封印）",
	RiftWaypointRetrieveAllError:          "rift：取得航點清單失敗：%s",
	RiftWaypointDestroySuccess:            "rift：航點 %q 已被銷毀",
	RiftWaypointDestroyError:              "rift：銷毀航點 %q 失敗，[ERROR: %s]",
	RiftWaypointRebindNotDirError:         "rift：重新綁定的路徑 %q 不是一個目錄",
	RiftWaypointRebindSuccess:             "rift：航點 %q 已重新綁定至 %s",
	RiftWaypointReforgeEmptyError:         "rift：重鑄名稱不能為空",
	RiftWaypointReforgeAlreadyExistsError: "rift：航點 %q 已存在，無法重鑄為已有名稱",
	RiftWaypointReforgeError:              "rift：重鑄航點 %q 失敗，[ERROR: %s]",
	RiftWaypointReforgeSuccess:            "rift：航點 %q 已重鑄為 %q",

	// Rune operations
	RiftRuneEngraveSuccessful: "rift：符文已刻印至航點 %q",
	RiftRuneEngraveNone:       "rift：未將符文刻印至航點 %q",
	RiftRuneUpdateError:       "rift：更新路徑 %q 的符文失敗，[ERROR: %s]",

	// Spell detail view
	RiftSpellDetailName:      "咒語名稱：",
	RiftSpellDetailCommand:   "咒語命令：",
	RiftSpellDetailAddedAt:   "咒語新增時間：",
	RiftSpellDetailCastCount: "咒語施放次數：",

	// Waypoint detail view
	RiftWaypointDetailName:           "航點名稱：",
	RiftWaypointDetailPath:           "航點路徑：",
	RiftWaypointDetailDiscovered:     "航點發現時間：",
	RiftWaypointDetailTravelledCount: "航點移動次數：",
	RiftWaypointDetailSealed:         "航點封印：",
	RiftWaypointDetailSealedReason:   "封印原因：",
	RiftWaypointDetailSealedTrue:     "是",
	RiftWaypointDetailSealedFalse:    "否",

	// Waypoint interactive UI
	WaypointInfoListTitle:                         "航點列表",
	WaypointInteractiveError:                      "rift：航點互動介面啟動失敗，[ERROR: %s]",
	RebindPathInputPlaceHolder:                    "輸入絕對路徑，留空則使用目前工作目錄",
	WaypointRebindTitle:                           "重新綁定",
	ReforgeWaypointNameInputPlaceHolder:           "輸入航點的新名稱",
	WaypointReforgeTitle:                          "重鑄",
	WaypointUIUpKeyHelp:                           "上移",
	WaypointUIUpKeyHelpDescription:                "將游標移至上一個航點",
	WaypointUIDownKeyHelp:                         "下移",
	WaypointUIDownKeyHelpDescription:              "將游標移至下一個航點",
	WaypointUIQuitKeyHelp:                         "退出",
	WaypointUIQuitKeyHelpDescription:              "退出航點互動介面",
	WaypointUIHelpKeyHelp:                         "說明",
	WaypointUIHelpKeyHelpDescription:              "顯示或隱藏完整的快捷鍵列表",
	WaypointNavigateKeyHelp:                       "導航",
	WaypointNavigateKeyHelpDescription:            "前往所選航點",
	WaypointDestroyKeyHelp:                        "銷毀",
	WaypointDestroyKeyHelpDescription:             "永久刪除所選航點",
	WaypointUnsealKeyHelp:                         "解封",
	WaypointUnsealKeyHelpDescription:              "嘗試解封所選航點",
	WaypointRebindKeyHelp:                         "重新綁定航點路徑",
	WaypointRebindKeyHelpDescription:              "將所選航點重新指定至新路徑",
	WaypointReforgeKeyHelp:                        "重鑄航點名稱",
	WaypointReforgeKeyHelpDescription:             "將所選航點重新命名為新名稱",
	WaypointNameCopyPathCopyKeyHelp:               "複製航點名稱 / 路徑至剪貼簿",
	WaypointNameCopyPathCopyKeyHelpDescription:    "將航點名稱 (y) 或路徑 (Y) 複製至剪貼簿",
	WaypointCopyFromInputValueKeyHelp:             "複製輸入內容",
	WaypointCopyFromInputValueKeyHelpDescription:  "將當前輸入框的內容複製至剪貼簿",
	WaypointPasteIntoInputValueKeyHelp:            "貼上至輸入框",
	WaypointPasteIntoInputValueKeyHelpDescription: "將剪貼簿內容貼上至輸入框",
	WaypointClosePopUp:                            "關閉彈窗",
	WaypointClosePopUpDescription:                 "關閉目前彈窗而不儲存",

	// Spell interactive UI
	SpellInfoListTitle:                       "法術書",
	SpellbookInteractiveError:                "rift：法術書互動介面啟動失敗，[ERROR: %s]",
	RiftSpellRetrieveAllError:                "rift：取得法術失敗",
	SpellUIUpKeyHelp:                         "上移",
	SpellUIUpKeyHelpDescription:              "將游標移至上一個法術",
	SpellUIDownKeyHelp:                       "下移",
	SpellUIDownKeyHelpDescription:            "將游標移至下一個法術",
	SpellUIQuitKeyHelp:                       "退出",
	SpellUIQuitKeyHelpDescription:            "退出法術書互動介面",
	SpellUIHelpKeyHelp:                       "說明",
	SpellUIHelpKeyHelpDescription:            "顯示或隱藏完整的快捷鍵列表",
	SpellCastKeyHelp:                         "施放",
	SpellCastKeyHelpDescription:              "在目前工作目錄施放法術，或施放至已發現的航點",
	SpellLearnKeyHelp:                        "習得",
	SpellLearnKeyHelpDescription:             "習得一個新法術",
	SpellForgetKey:                           "遺忘",
	SpellForgetKeyDescription:                "永久遺忘所選法術",
	SpellClosePopUp:                          "關閉彈窗",
	SpellClosePopUpDescription:               "關閉目前彈窗而不儲存",
	SpellUILearnKeyHelp:                      "確認",
	SpellUINextInputKeyHelp:                  "下一個輸入",
	SpellUIPreviousInputKeyHelp:              "上一個輸入",
	SpellNameInputPlaceHolder:                "輸入咒語的名稱",
	SpellCommandInputPlaceHolder:             "輸入咒語的命令",
	SpellNameInputTitle:                      "咒語名稱：",
	SpellCommandInputTitle:                   "咒語命令：",
	SpellUIChooseCastLocationKeyHelp:         "選擇施放位置",
	SpellUIChooseWaypointCastLocationKeyHelp: "選擇施放咒語的航點",

	// Rune interactive UI
	RuneInteractiveError:           "[ERROR: %s]",
	RuneEngraveTypeOptionListTitle: "符文選項",
	EngraveRuneEnterTitle:          "進入時符文的命令：",
	EngraveRuneLeaveTitle:          "離開時符文的命令：",
	EngraveRuneEngraveButton:       "刻印",
	RuneCommandsPlaceHolder:        "輸入命令…（cd 無效，推薦使用 rift 切換路徑）",
	RuneCommandsInvalidDueToShellBuildInCommand: "偵測到符文使用了 Shell 內建命令（如 cd、export、source、alias）——內建命令只影響執行它們的進程，無法修改當前 Shell 會話。" +
		"如需在命令序列中使用內建命令，請使用 Shell 的 -c 參數顯式呼叫 Shell 並串連命令。" +
		"使用 -i（互動模式）可載入 Shell 的互動設定（.zshrc、.bashrc 等）——若環境依賴其中內容（如使用 rift 切換目錄、nvm、conda）則必須使用此選項。\n\n" +
		"範例:\n" +
		"  bash -i -c \"source env/bin/activate && python main.py\"\n" +
		"  zsh -i -c \"source env/bin/activate && python main.py\"\n" +
		"  fish -i -c \"source env/bin/activate.fish; python main.py\"\n" +
		"  pwsh -Login -Command \". ./env/bin/Activate.ps1; python main.py\"\n" +
		"  cmd /c \"activate.bat && python main.py\"",
	EngraveRuneEnterOptionName: "刻印進入時符文",
	EngraveRuneEnterOptionDesc: "設定進入此航點時執行的命令",
	EngraveRuneLeaveOptionName: "刻印離開時符文",
	EngraveRuneLeaveOptionDesc: "設定離開此航點時執行的命令",
	RemoveRuneEnterOptionName:  "移除進入時符文",
	RemoveRuneEnterOptionDesc:  "清除進入此航點時執行的命令",
	RemoveRuneLeaveOptionName:  "移除離開時符文",
	RemoveRuneLeaveOptionDesc:  "清除離開此航點時執行的命令",

	// Cast location option popup
	CastLocationOptionTitle:               "施放位置",
	CastLocationOptionCurrent:             "當前目錄",
	CastLocationOptionCurrentDescription:  "在當前工作目錄中施放咒語",
	CastLocationOptionWaypoint:            "航點",
	CastLocationOptionWaypointDescription: "從已發現的航點中選擇並在那裡施放咒語",
	CastWaypointLocationOptionTitle:       "選擇航點",

	// Setup
	CheckAndRunSetupError:  "rift：設定失敗，[ERROR: %s]",
	RiftAutoSetupTriggered: "rift：設定與配置已自動觸發",
}
