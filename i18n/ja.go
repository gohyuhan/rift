package i18n

var jA = LanguageMapping{
	// General
	ConfigPathError:          "設定パスが無効です、[ERROR: %s]",
	RiftReservedKeywordError: "`%s` は rift の予約済みキーワードです",
	RiftDetectedShell:        "rift：検出されたシェル：%s",

	// Settings related
	SettingsPathError:  "設定ディレクトリへのアクセスに失敗しました、[ERROR: %s]",
	SettingsReadError:  "設定ファイルの読み込みに失敗しました、[ERROR: %s]",
	SettingsParseError: "設定ファイルの解析に失敗しました。デフォルトにリセットします、[ERROR: %s]",
	SettingsWriteError: "設定ファイルの書き込みに失敗しました、[ERROR: %s]",

	// DB related
	DBPathError:                 "データベースディレクトリへのアクセスに失敗しました、[ERROR: %s]",
	DBSetupError:                "データベースの初期化に失敗しました、[ERROR: %s]",
	DBOpenError:                 "データベースのオープンに失敗しました。初期化されていない可能性があります。`rift awaken` を実行してください",
	SettingsBucketNotFoundError: "データベースに設定バケットが見つかりません。`rift awaken` を再実行してください",
	WaypointBucketNotFoundError: "データベースにウェイポイントバケットが見つかりません。`rift awaken` を再実行してください",

	// Updater related
	UpdaterDownloadPrompt:               "新しいバージョン %s が利用可能です。今すぐダウンロードしますか？(y/n): ",
	UpdaterFailToCheckForUpdate:         "アップデートの確認に失敗しました：%v",
	UpdaterFailToCreateRequest:          "リクエストの作成に失敗しました：%v",
	UpdaterFailToFetchRelease:           "最新リリース情報の取得に失敗しました：%v",
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
	RiftDescription:       "事前に定義したチェックポイント名でパスを簡単に移動できます",
	RiftAwakenDescription: "シェル内で rift を起動します【初回使用時のセットアップと初期化を行います】",

	// cmd root
	RiftSavedWaypoint:   "rift：%q -> %s を保存しました",
	RiftUnknownWaypoint: "rift：不明なウェイポイント名 %q",

	// Setup related
	CheckAndRunSetupError: "rift：セットアップに失敗しました、[ERROR: %s]",
}
