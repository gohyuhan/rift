package updater

// ----------------------------------
//
//	Checks whether an update is due, fetches the latest release if so,
//	and prompts the user to download it when a newer version is available.
//
// ----------------------------------
func AutoUpdater() {
	if ShouldCheckForUpdate() {
		latestVersion, isNewer, err := CheckForUpdates()
		if err != nil {
			return
		}

		if isNewer {
			if PromptUserForUpdate(latestVersion) {
				Update()
			}
		}
	}
}
