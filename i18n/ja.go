package i18n

var jA = LanguageMapping{
	// General
	ConfigPathError:          "設定パスが無効です、[ERROR: %s]",
	RiftReservedKeywordError: "`%s` は rift の予約済みキーワードです",
	RiftDetectedShell:        "rift：検出されたシェル：%s",
	CWDIsNotDirError:         "現在の作業ディレクトリは有効なディレクトリではありません",
	PathNotAbsoluteError:     "パスは絶対パスである必要があります。指定されたパス: %s",
	NotFileOrDirError:        "指定されたパスはファイルまたはディレクトリとして存在しません",

	// Settings related
	SettingsPathError:                 "設定ディレクトリへのアクセスに失敗しました、[ERROR: %s]",
	SettingsReadError:                 "設定ファイルの読み込みに失敗しました、[ERROR: %s]",
	SettingsParseError:                "設定ファイルの解析に失敗しました。デフォルトにリセットします、[ERROR: %s]",
	SettingsWriteError:                "設定ファイルの書き込みに失敗しました、[ERROR: %s]",
	SettingsLanguageUpdated:           "rift: 言語を %s に設定しました",
	SettingsLanguageNotSupported:      "rift: 言語 %q はサポートされていません（対応言語: %s）",
	SettingsAutoUpdateUpdated:         "rift: 自動アップデートを %t に設定しました",
	SettingsDownloadPreReleaseUpdated: "rift: プレリリースのダウンロードを %t に設定しました",

	// DB related
	DBPathError:                 "データベースディレクトリへのアクセスに失敗しました、[ERROR: %s]",
	DBSetupError:                "データベースの初期化に失敗しました、[ERROR: %s]",
	DBOpenError:                 "データベースのオープンに失敗しました。初期化されていない可能性があります。`rift awaken` を実行してください",
	SettingsBucketNotFoundError: "データベースに設定バケットが見つかりません。`rift awaken` を再実行してください",
	WaypointBucketNotFoundError: "データベースにウェイポイントバケットが見つかりません。`rift awaken` を再実行してください",
	WaypointDataCorruptedError:  "ウェイポイント %q のデータが破損しており、読み込めません",

	// Updater related
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

	// Shell related
	ShellCMDNotSupported:  "Windows コマンドプロンプトはシェル関数をサポートしていません。\nPowerShell、Git Bash、または WSL を使用し、`rift awaken` を再実行してください。",
	ShellUnsupported:      "シェル %q は rift でサポートされていません。\nサポートされているシェル：bash、zsh、fish、ksh、PowerShell。\n手動で統合を追加できます — docs/shell-integration.md を参照してください",
	ShellNoConfigFile:     "シェル %q の既知の設定ファイルがありません",
	ShellAlreadyInstalled: "rift：シェル統合はすでに %s に存在します",
	ShellInstallSuccess:   "rift：シェル統合を %s に追加しました",
	ShellInstallReload:    "rift：シェルを再起動するか、次を実行してください：%s",
	BinaryNotInPath:       "rift：PATH にバイナリが見つかりません — このセッション後も使用するには rift を PATH に追加してください",

	// cmd description
	RiftDescription:                       "事前に定義したチェックポイント名でパスを簡単に移動できます",
	RiftAwakenDescription:                 "シェル内で rift を起動します【初回使用時のセットアップと初期化を行います】",
	RiftDiscoverDescription:               "現在の作業ディレクトリにウェイポイント名を割り当てます",
	RiftWaypointDescription:               "すべてのウェイポイントを一覧表示するか、特定のウェイポイントの情報を表示します",
	RiftFlagLanguageDescription:           "rift の言語を設定します（対応言語：EN、JA、ZH-HANS、ZH-HANT）",
	RiftFlagAutoUpdateDescription:         "rift が自動的にアップデートを確認するかどうかを設定します（有効にするには --autoupdate、無効にするには --autoupdate=false）",
	RiftFlagDownloadPreReleaseDescription: "rift がプレリリース版もダウンロードするか、安定版のみにするかを設定します（有効にするには --download-pre-release、無効にするには --download-pre-release=false）",
	RiftFlagWaypointDestroyDescription:    "名前を指定してウェイポイントを削除します",
	RiftFlagWaypointRebindDescription:     "既存のウェイポイントを新しいパスに再割り当てします。パスが指定されない場合は現在の作業ディレクトリを使用し、有効な絶対パスが指定された場合はそのパスを優先します",
	RiftFlagWaypointReforgeDescription:    "既存のウェイポイントを新しい名前に変更します",
	RiftFlagUpdateDescription:             "最新バージョンの確認を手動でトリガーし、利用可能な場合はアップデートします",
	RiftFlagVersionDescription:            "rift の現在のバージョンを表示します",

	// Waypoint related
	RiftSavedWaypoint:                     "rift：%q -> %s を保存しました",
	RiftUnknownWaypoint:                   "rift：不明なウェイポイント名 %q",
	RiftWaypointAlreadyExistsError:        "ウェイポイント %q は既に存在し、%s を指しています",
	RiftWaypointDoNotExistsError:          "rift：ウェイポイント %q は存在しません",
	RiftWaypointUpdateError:               "rift：ウェイポイント %q の更新に失敗しました",
	RiftWaypointSealedError:               "rift：ウェイポイント %q は封印されており、移動できません。理由：%q",
	RiftWaypointSealedLabel:               "(封印済み)",
	RiftWaypointRetrieveAllError:          "rift：ウェイポイントの取得に失敗しました",
	RiftWaypointDetailName:                "ウェイポイント名：",
	RiftWaypointDetailPath:                "ウェイポイントパス：",
	RiftWaypointDetailDiscovered:          "ウェイポイント発見日時：",
	RiftWaypointDetailTravelledCount:      "ウェイポイント移動回数：",
	RiftWaypointDetailSealed:              "ウェイポイント封印：",
	RiftWaypointDetailSealedReason:        "封印理由：",
	RiftWaypointDetailSealedTrue:          "はい",
	RiftWaypointDetailSealedFalse:         "いいえ",
	RiftWaypointDestroySuccess:            "rift：ウェイポイント %q を削除しました",
	RiftWaypointDestroyError:              "rift：ウェイポイント %q の削除に失敗しました、[ERROR: %s]",
	RiftWaypointRebindNotDirError:         "rift：再バインド先のパス %q はディレクトリではありません",
	RiftWaypointRebindSuccess:             "rift：ウェイポイント %q を %s に再バインドしました",
	RiftWaypointRebindError:               "rift：ウェイポイント %q の再バインドに失敗しました、[ERROR: %s]",
	RiftWaypointReforgeEmptyError:         "rift：リフォージ名を空にすることはできません",
	RiftWaypointReforgeError:              "rift：ウェイポイント %q のリフォージに失敗しました、[ERROR: %s]",
	RiftWaypointReforgeAlreadyExistsError: "rift：ウェイポイント %q は既に存在するため、既存の名前にリフォージすることはできません",
	RiftWaypointReforgeSuccess:            "rift：ウェイポイント %q を %q にリフォージしました",

	// Flag related
	RiftFlagRetrieveError: "rift：フラグ %q の取得に失敗しました、[ERROR: %s]",

	// Setup related
	CheckAndRunSetupError:  "rift：セットアップに失敗しました、[ERROR: %s]",
	RiftAutoSetupTriggered: "rift：設定とコンフィグのセットアップが自動的にトリガーされました",
}
