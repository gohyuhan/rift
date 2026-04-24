package i18n

var jA = LanguageMapping{
	// General
	ConfigPathError:          "設定パスが無効です、[ERROR: %s]",
	RiftReservedKeywordError: "`%s` は rift の予約済みキーワードです",
	RiftDetectedShell:        "rift：検出されたシェル：%s",
	CWDIsNotDirError:         "現在の作業ディレクトリは有効なディレクトリではありません",
	PathNotAbsoluteError:     "パスは絶対パスである必要があります。指定されたパス: %s",
	NotFileOrDirError:        "指定されたパスはファイルまたはディレクトリとして存在しません",
	InvalidValueProvided:     "無効な値が指定されました。スペースは使用できません。また、空にすることもできません",
	SkippingDueToExecutorErr: "rift：エグゼキューターの起動に失敗したため、ルーンコマンドをスキップします",

	// Settings
	SettingsPathError:                 "設定ディレクトリへのアクセスに失敗しました、[ERROR: %s]",
	SettingsReadError:                 "設定ファイルの読み込みに失敗しました、[ERROR: %s]",
	SettingsParseError:                "設定ファイルの解析に失敗しました。デフォルトにリセットします、[ERROR: %s]",
	SettingsWriteError:                "設定ファイルの書き込みに失敗しました、[ERROR: %s]",
	SettingsLanguageUpdated:           "rift: 言語を %s に設定しました",
	SettingsLanguageNotSupported:      "rift: 言語 %q はサポートされていません（対応言語: %s）",
	SettingsAutoUpdateUpdated:         "rift: 自動アップデートを %t に設定しました",
	SettingsDownloadPreReleaseUpdated: "rift: プレリリースのダウンロードを %t に設定しました",

	// Database
	DBPathError:                 "データベースディレクトリへのアクセスに失敗しました、[ERROR: %s]",
	DBSetupError:                "データベースの初期化に失敗しました、[ERROR: %s]",
	DBOpenError:                 "データベースのオープンに失敗しました — 以前のセッションが正常に終了せずロックが残っているか、別の rift プロセスがすでに起動中の可能性があります。rift がクラッシュまたは強制終了された場合は、`lsof | grep rift.db` を実行してプロセスを特定・終了してから再試行してください。初めてのインストールの場合は `rift awaken` を実行してください",
	WaypointBucketNotFoundError: "データベースにウェイポイントバケットが見つかりません。`rift awaken` を再実行してください",
	WaypointDataCorruptedError:  "ウェイポイント [%s] のデータが破損しており、読み込めません",
	SpellBucketNotFoundError:    "データベースにスペルバケットが見つかりません。`rift awaken` を再実行してください",
	SpellDataCorruptedError:     "スペル [%s] のデータが破損しており、読み込めません",
	RuneBucketNotFoundError:     "データベースにルーンバケットが見つかりません。`rift awaken` を再実行してください",
	RuneDataCorruptedError:      "パス [%s] のルーンデータが破損しており、読み込めません",

	// Updater
	UpdaterDownloadPrompt:               "新しいバージョン %s が利用可能です。今すぐダウンロードしますか？(y/n): ",
	UpdaterFailToCheckForUpdate:         "アップデートの確認に失敗しました：%v",
	UpdaterFailToCreateRequest:          "リクエストの作成に失敗しました：%v",
	UpdaterFailToFetchRelease:           "最新リリース情報の取得に失敗しました：%v",
	UpdaterNoReleasesFound:              "リリースが見つかりません",
	UpdaterFailToReadResponse:           "レスポンスボディの読み取りに失敗しました：%v",
	UpdaterFailToParseJSON:              "JSON レスポンスの解析に失敗しました：%v",
	UpdaterFailToExtractBinary:          "バイナリの展開に失敗しました：%v",
	UpdaterUnsupportedArchiveFormat:     "サポートされていないアーカイブ形式です",
	UpdaterBinaryNotFoundInArchive:      "アーカイブ内にバイナリが見つかりません",
	UpdaterAlreadyLatest:                "すでに最新バージョン（%s）です",
	UpdaterDownloading:                  "バージョン %s をダウンロード中...",
	UpdaterUnSupportedOS:                "サポートされていないOS/アーキテクチャ：%s/%s",
	UpdaterDownloadFail:                 "アップデートのダウンロードに失敗しました：%v",
	UpdaterBinaryReplaceFail:            "バイナリの置き換えに失敗しました：%v",
	UpdaterDownloadSuccess:              "バージョン %s へのアップデートが完了しました",
	UpdaterDownloadUnexpectedStatusCode: "予期しないステータスコード：%d",
	UpdaterRequiresSudo:                 "権限がありません。sudoで再試行します...",

	// Shell
	ShellCMDNotSupported:  "Windows コマンドプロンプトはシェル関数をサポートしていません。\nPowerShell、Git Bash、または WSL を使用し、`rift awaken` を再実行してください。",
	ShellUnsupported:      "シェル %q は rift でサポートされていません。\nサポートされているシェル：bash、zsh、fish、ksh、PowerShell。\n手動で統合を追加できます — docs/shell-integration.md を参照してください",
	ShellNoConfigFile:     "シェル %q の既知の設定ファイルがありません",
	ShellAlreadyInstalled: "rift：シェル統合はすでに %s に存在します",
	ShellInstallSuccess:   "rift：シェル統合を %s に追加しました",
	ShellReloadHint:       "rift：シェルを再起動するか、次を実行してください：%s",
	BinaryNotInPath:       "rift：PATH にバイナリが見つかりません — このセッション後も使用するには rift を PATH に追加してください",

	// Commands and flags
	RiftDescription:                       "事前に定義したチェックポイント名でパスを簡単に移動できます",
	RiftAwakenDescription:                 "シェル内で rift を起動します【初回使用時のセットアップと初期化を行います】",
	RiftDiscoverDescription:               "現在の作業ディレクトリにウェイポイント名を割り当てます",
	RiftWaypointDescription:               "ウェイポイントのインタラクティブUIを起動するか、特定のウェイポイントの情報を表示します",
	RiftLearnDescription:                  "コマンドに名前を付けて rift に新しいスペルを教えます。複数単語のコマンドは引用符で囲んでください（例：rift learn build \"docker compose up --build\"）",
	RiftSpellDescription:                  "スペル名を指定してキャストし、紐付けられたターミナルコマンドを実行します",
	RiftSpellbookDescription:              "スペルのインタラクティブUIを起動するか、特定のスペルの情報を表示します",
	RiftFlagLanguageDescription:           "rift の言語を設定します（対応言語：EN、JA、ZH-HANS、ZH-HANT）",
	RiftFlagAutoUpdateDescription:         "rift が自動的にアップデートを確認するかどうかを設定します（有効にするには --autoupdate、無効にするには --autoupdate=false）",
	RiftFlagDownloadPreReleaseDescription: "rift がプレリリース版もダウンロードするか、安定版のみにするかを設定します（有効にするには --download-pre-release、無効にするには --download-pre-release=false）",
	RiftFlagWaypointDestroyDescription:    "名前を指定してウェイポイントを削除します",
	RiftFlagSpellForgetDescription:        "名前を指定して保存済みスペルを削除します",
	RiftFlagWaypointRebindDescription:     "既存のウェイポイントを新しいパスに再割り当てします。パスが指定されない場合は現在の作業ディレクトリを使用し、有効な絶対パスが指定された場合はそのパスを優先します",
	RiftFlagWaypointReforgeDescription:    "既存のウェイポイントを新しい名前に変更します",
	RiftFlagUpdateDescription:             "最新バージョンの確認を手動でトリガーし、利用可能な場合はアップデートします",
	RiftFlagVersionDescription:            "rift の現在のバージョンを表示します",
	RiftFlagCastDescription:               "ナビゲーションの代わりに、習得したスペルをウェイポイントのパスでキャストします。スペルコマンドは、ウェイポイントのパスを作業ディレクトリとして実行されます",
	RiftFlagRetrieveError:                 "rift：フラグ %q の取得に失敗しました、[ERROR: %s]",
	RiftRuneDescription:                   "ウェイポイントに移動時・離脱時のトリガーコマンドを設定します；rift でそのウェイポイントへ、またはそこから移動する際に自動的に実行されます",

	// Spell operations
	RiftSpellSaved:            "rift：%q -> %s を習得しました",
	RiftSpellUpdated:          "rift：スペル %q を更新しました -> %s",
	RiftSpellForgetSuccess:    "rift：スペル %q を忘却しました",
	RiftSpellForgetError:      "rift：スペル %q の忘却に失敗しました、[ERROR: %s]",
	RiftSpellDoNotExistsError: "rift：スペル %q は存在しません",
	RiftSpellUpdateError:      "rift：スペル %q の更新に失敗しました、[ERROR: %s]",
	ForbiddenShellBuiltinSpellCommand: "rift：シェル組み込みコマンド（cd、export、source、alias など）は実行したプロセス内にのみ影響し、現在のシェルセッションを変更することはできません。" +
		"組み込みコマンドをコマンドシーケンスの一部として使用するには、シェルを -c フラグで明示的に呼び出してコマンドをチェーンしてください。" +
		"--login（または同等のオプション）を使用すると、完全なシェル環境（PATH、エイリアス、プロファイルなど）を読み込めます。\n\n" +
		"例:\n" +
		"  bash --login -c \"source env/bin/activate && python main.py\"\n" +
		"  zsh --login -c \"source env/bin/activate && python main.py\"\n" +
		"  fish --login -c \"source env/bin/activate.fish; python main.py\"\n" +
		"  pwsh -Login -Command \". ./env/bin/Activate.ps1; python main.py\"\n" +
		"  cmd /c \"activate.bat && python main.py\"",
	ForbiddenRiftNavigationSpellCommand: "rift：rift ウェイポイントナビゲーションコマンド（例: rift <ウェイポイント名>）はスペルとして登録できません — これらはシェル統合によって直接処理され、子プロセスとして実行されません。",
	ForbiddenRiftNavigationRuneCommand:  "rift：rift ウェイポイントナビゲーションコマンド（例: rift <ウェイポイント名>）はルーンではサポートされていません — 特定のディレクトリでコマンドを実行するには、`rift <ウェイポイント名> --cast <スペル名>` を使用してください。",
	SpellCommandEmpty:                   "rift：スペルコマンドは空にできません",
	InvalidSpellCommandError:            "rift：スペルコマンド [%s] が無効のため実行できませんでした",

	// Waypoint operations
	RiftSavedWaypoint:                     "rift：%q -> %s を保存しました",
	RiftWaypointAlreadyExistsError:        "ウェイポイント %q は既に存在し、%s を指しています",
	RiftWaypointDoNotExistsError:          "rift：ウェイポイント %q は存在しません",
	RiftWaypointUpdateError:               "rift：ウェイポイント %q の更新に失敗しました、[ERROR: %s]",
	RiftWaypointSealedError:               "rift：ウェイポイント %q は封印されており、移動できません。理由：%q",
	RiftWaypointSealedLabel:               "(封印済み)",
	RiftWaypointRetrieveAllError:          "rift：ウェイポイントの取得に失敗しました: %s",
	RiftWaypointDestroySuccess:            "rift：ウェイポイント %q を削除しました",
	RiftWaypointDestroyError:              "rift：ウェイポイント %q の削除に失敗しました、[ERROR: %s]",
	RiftWaypointRebindNotDirError:         "rift：再バインド先のパス %q はディレクトリではありません",
	RiftWaypointRebindSuccess:             "rift：ウェイポイント %q を %s に再バインドしました",
	RiftWaypointReforgeEmptyError:         "rift：リフォージ名を空にすることはできません",
	RiftWaypointReforgeAlreadyExistsError: "rift：ウェイポイント %q は既に存在するため、既存の名前にリフォージすることはできません",
	RiftWaypointReforgeError:              "rift：ウェイポイント %q のリフォージに失敗しました、[ERROR: %s]",
	RiftWaypointReforgeSuccess:            "rift：ウェイポイント %q を %q にリフォージしました",

	// Rune operations
	RiftRuneEngraveSuccessful: "rift：ルーンをウェイポイント %q に刻みました",
	RiftRuneEngraveNone:       "rift：ウェイポイント %q にルーンは刻まれませんでした",
	RiftRuneUpdateError:       "rift：パス %q のルーン更新に失敗しました、[ERROR: %s]",

	// Spell detail view
	RiftSpellDetailName:      "スペル名：",
	RiftSpellDetailCommand:   "スペルコマンド：",
	RiftSpellDetailAddedAt:   "スペル追加日時：",
	RiftSpellDetailCastCount: "スペルキャスト回数：",

	// Waypoint detail view
	RiftWaypointDetailName:           "ウェイポイント名：",
	RiftWaypointDetailPath:           "ウェイポイントパス：",
	RiftWaypointDetailDiscovered:     "ウェイポイント発見日時：",
	RiftWaypointDetailTravelledCount: "ウェイポイント移動回数：",
	RiftWaypointDetailSealed:         "ウェイポイント封印：",
	RiftWaypointDetailSealedReason:   "封印理由：",
	RiftWaypointDetailSealedTrue:     "はい",
	RiftWaypointDetailSealedFalse:    "いいえ",

	// Waypoint interactive UI
	WaypointInfoListTitle:                         "ウェイポイント一覧",
	WaypointInteractiveError:                      "rift：ウェイポイントインタラクティブセッションの起動に失敗しました、[ERROR: %s]",
	RebindPathInputPlaceHolder:                    "絶対パスを入力してください。空のままにすると現在の作業ディレクトリが使用されます",
	WaypointRebindTitle:                           "リバインド",
	ReforgeWaypointNameInputPlaceHolder:           "ウェイポイントの新しい名前を入力してください",
	WaypointReforgeTitle:                          "リフォージ",
	WaypointUIUpKeyHelp:                           "上へ",
	WaypointUIUpKeyHelpDescription:                "前のウェイポイントにカーソルを移動する",
	WaypointUIDownKeyHelp:                         "下へ",
	WaypointUIDownKeyHelpDescription:              "次のウェイポイントにカーソルを移動する",
	WaypointUIQuitKeyHelp:                         "終了",
	WaypointUIQuitKeyHelpDescription:              "ウェイポイントインタラクティブUIを終了する",
	WaypointUIHelpKeyHelp:                         "ヘルプ",
	WaypointUIHelpKeyHelpDescription:              "キーバインド一覧の表示・非表示を切り替える",
	WaypointNavigateKeyHelp:                       "移動",
	WaypointNavigateKeyHelpDescription:            "選択したウェイポイントへ移動する",
	WaypointDestroyKeyHelp:                        "破壊",
	WaypointDestroyKeyHelpDescription:             "選択したウェイポイントを完全に削除する",
	WaypointUnsealKeyHelp:                         "封印解除",
	WaypointUnsealKeyHelpDescription:              "選択したウェイポイントの封印を解除する",
	WaypointRebindKeyHelp:                         "パスをリバインド",
	WaypointRebindKeyHelpDescription:              "選択したウェイポイントを新しいパスに再割り当てする",
	WaypointReforgeKeyHelp:                        "名前をリフォージ",
	WaypointReforgeKeyHelpDescription:             "選択したウェイポイントを新しい名前に変更する",
	WaypointNameCopyPathCopyKeyHelp:               "ウェイポイント名 / パスをクリップボードにコピー",
	WaypointNameCopyPathCopyKeyHelpDescription:    "ウェイポイント名 (y) またはパス (Y) をクリップボードにコピーする",
	WaypointCopyFromInputValueKeyHelp:             "入力値をコピー",
	WaypointCopyFromInputValueKeyHelpDescription:  "現在の入力フィールドの内容をクリップボードにコピーする",
	WaypointPasteIntoInputValueKeyHelp:            "入力欄に貼り付け",
	WaypointPasteIntoInputValueKeyHelpDescription: "クリップボードの内容を入力フィールドに貼り付ける",
	WaypointClosePopUp:                            "ポップアップを閉じる",
	WaypointClosePopUpDescription:                 "保存せずに現在のポップアップを閉じる",

	// Spell interactive UI
	SpellInfoListTitle:                       "魔法書",
	SpellbookInteractiveError:                "rift：スペルブックインタラクティブセッションの起動に失敗しました、[ERROR: %s]",
	RiftSpellRetrieveAllError:                "rift：スペルの取得に失敗しました",
	SpellUIUpKeyHelp:                         "上へ",
	SpellUIUpKeyHelpDescription:              "前のスペルにカーソルを移動する",
	SpellUIDownKeyHelp:                       "下へ",
	SpellUIDownKeyHelpDescription:            "次のスペルにカーソルを移動する",
	SpellUIQuitKeyHelp:                       "終了",
	SpellUIQuitKeyHelpDescription:            "スペルブックインタラクティブUIを終了する",
	SpellUIHelpKeyHelp:                       "ヘルプ",
	SpellUIHelpKeyHelpDescription:            "キーバインド一覧の表示・非表示を切り替える",
	SpellCastKeyHelp:                         "キャスト",
	SpellCastKeyHelpDescription:              "現在の作業ディレクトリでスペルをキャストする、またはウェイポイントへキャストする",
	SpellLearnKeyHelp:                        "習得",
	SpellLearnKeyHelpDescription:             "新しいスペルを習得する",
	SpellForgetKey:                           "忘却",
	SpellForgetKeyDescription:                "選択したスペルを完全に忘却する",
	SpellClosePopUp:                          "ポップアップを閉じる",
	SpellClosePopUpDescription:               "保存せずに現在のポップアップを閉じる",
	SpellUILearnKeyHelp:                      "確定",
	SpellUINextInputKeyHelp:                  "次の入力",
	SpellUIPreviousInputKeyHelp:              "前の入力",
	SpellNameInputPlaceHolder:                "スペルの名前を入力してください",
	SpellCommandInputPlaceHolder:             "スペルのコマンドを入力してください",
	SpellNameInputTitle:                      "スペル名：",
	SpellCommandInputTitle:                   "スペルコマンド：",
	SpellUIChooseCastLocationKeyHelp:         "詠唱場所を選択",
	SpellUIChooseWaypointCastLocationKeyHelp: "スペルを詠唱するウェイポイントを選択",

	// Rune interactive UI
	RuneInteractiveError:           "[ERROR: %s]",
	RuneEngraveTypeOptionListTitle: "ルーンオプション",
	EngraveRuneEnterTitle:          "移動時ルーンのコマンド：",
	EngraveRuneLeaveTitle:          "離脱時ルーンのコマンド：",
	EngraveRuneEngraveButton:       "刻む",
	RuneCommandsPlaceHolder:        "コマンドを入力… （cd は効果がありません。パス変更には rift の使用を推薦します）",
	RuneCommandsInvalidDueToShellBuildInCommand: "シェル組み込みコマンド（cd、export、source、alias など）が検出されました。組み込みコマンドは実行したプロセス内にのみ影響し、現在のシェルセッションを変更することはできません。" +
		"組み込みコマンドをコマンドシーケンスの一部として使用するには、シェルを -c フラグで明示的に呼び出してコマンドをチェーンしてください。" +
		"-i（インタラクティブモード）を使用すると、シェルのインタラクティブ設定（.zshrc、.bashrc など）を読み込めます — 環境がそれに依存している場合（rift でのディレクトリ移動、nvm、conda など）は必須です。\n\n" +
		"例:\n" +
		"  bash -i -c \"source env/bin/activate && python main.py\"\n" +
		"  zsh -i -c \"source env/bin/activate && python main.py\"\n" +
		"  fish -i -c \"source env/bin/activate.fish; python main.py\"\n" +
		"  pwsh -Login -Command \". ./env/bin/Activate.ps1; python main.py\"\n" +
		"  cmd /c \"activate.bat && python main.py\"",
	EngraveRuneEnterOptionName: "移動時ルーンを刻む",
	EngraveRuneEnterOptionDesc: "このウェイポイントに移動した際に実行するコマンドを設定する",
	EngraveRuneLeaveOptionName: "離脱時ルーンを刻む",
	EngraveRuneLeaveOptionDesc: "このウェイポイントから離脱した際に実行するコマンドを設定する",
	RemoveRuneEnterOptionName:  "移動時ルーンを削除する",
	RemoveRuneEnterOptionDesc:  "このウェイポイントに移動した際に実行するコマンドを削除する",
	RemoveRuneLeaveOptionName:  "離脱時ルーンを削除する",
	RemoveRuneLeaveOptionDesc:  "このウェイポイントから離脱した際に実行するコマンドを削除する",

	// Cast location option popup
	CastLocationOptionTitle:               "詠唱場所",
	CastLocationOptionCurrent:             "現在のディレクトリ",
	CastLocationOptionCurrentDescription:  "現在の作業ディレクトリでスペルを詠唱する",
	CastLocationOptionWaypoint:            "ウェイポイント",
	CastLocationOptionWaypointDescription: "発見済みのウェイポイントから選択し、そこでスペルを詠唱する",
	CastWaypointLocationOptionTitle:       "ウェイポイントを選択",

	// Setup
	CheckAndRunSetupError:  "rift：セットアップに失敗しました、[ERROR: %s]",
	RiftAutoSetupTriggered: "rift：設定とコンフィグのセットアップが自動的にトリガーされました",
}
