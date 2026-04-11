package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
)

// ----------------------------------
//
//	builds or rebuilds the bubbles list model on the WaypointInteractiveModel;
//	preserves the previously selected waypoint by name when re-initialising,
//	and falls back to the stored cursor position when no match is found;
//	called once after the first WindowSizeMsg so that layout dimensions
//	are valid before the list is rendered
//
// ----------------------------------
func initWaypointInfoListModel(m *WaypointInteractiveModel) error {
	previousSelectedWaypoint := m.WaypointInfoList.SelectedItem()
	selectedWayPointCursorPosition := -1

	latestWaypointInfoArray := []list.Item{}

	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - ListTitleHorizontalPadding

	allWaypointsInfo, err := apiUtils.GetAllWaypointsInfo(m.BboltReadDb)
	if err != nil {
		return err
	}

	if previousSelectedWaypoint != nil {
		// a previous selection exists; scan the fresh list to find a name match
		// so the cursor lands on the same waypoint after re-initialisation
		previousSelectedWaypointInfo := previousSelectedWaypoint.(waypointInfoItem)
		for index, waypoint := range allWaypointsInfo {
			if waypoint.WaypointName == previousSelectedWaypointInfo.WaypointName {
				selectedWayPointCursorPosition = index
			}
			latestWaypointInfoArray = append(latestWaypointInfoArray, waypointInfoItem(waypoint))
		}
	} else {
		// no previous selection; convert all records to list items
		for _, waypoint := range allWaypointsInfo {
			latestWaypointInfoArray = append(latestWaypointInfoArray, waypointInfoItem(waypoint))
		}
	}

	m.WaypointInfoList = list.New(latestWaypointInfoArray, waypointInfoDelegate{}, m.Width, m.Height)
	m.WaypointInfoList.SetShowPagination(false)
	m.WaypointInfoList.SetShowStatusBar(false)
	m.WaypointInfoList.SetFilteringEnabled(false)
	m.WaypointInfoList.SetShowFilter(false)
	m.WaypointInfoList.SetShowHelp(true)

	// truncate the title to prevent overflow when the terminal is narrow
	m.WaypointInfoList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.WaypointInfoListTitle, titleWidthLimit, "...")
	m.WaypointInfoList.Styles.Title = style.NewStyle.Bold(true)
	m.WaypointInfoList.Styles.PaginationStyle = style.NewStyle
	m.WaypointInfoList.Styles.TitleBar = style.NewStyle
	m.WaypointInfoList.Styles.HelpStyle = style.NewStyle
	m.WaypointInfoList.Help.Styles.ShortKey = style.NewStyle.Foreground(style.ColorBlueMuted).Bold(true)
	m.WaypointInfoList.Help.Styles.ShortDesc = style.NewStyle.Foreground(style.ColorBlueMuted)
	m.WaypointInfoList.Help.Styles.ShortSeparator = style.NewStyle.Foreground(style.ColorBlueGrayMuted)
	// clear the default key map so only our custom bindings appear in the help bar
	m.WaypointInfoList.KeyMap = list.KeyMap{}
	m.WaypointInfoList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
		}
	}

	if selectedWayPointCursorPosition >= 0 {
		// restore the cursor to the previously selected waypoint by name match
		m.WaypointInfoList.Select(selectedWayPointCursorPosition)
		m.WaypointInfoListCursorPosition = selectedWayPointCursorPosition
	} else {
		// clamp the stored position to the last item if the list has shrunk
		if m.WaypointInfoListCursorPosition > len(m.WaypointInfoList.Items())-1 {
			m.WaypointInfoList.Select(len(m.WaypointInfoList.Items()) - 1)
			m.WaypointInfoListCursorPosition = len(m.WaypointInfoList.Items()) - 1
		} else {
			m.WaypointInfoList.Select(m.WaypointInfoListCursorPosition)
		}
	}

	// update the help key map to reflect whether the initial selection is sealed
	currentSelectedWaypoint := m.WaypointInfoList.SelectedItem()
	if currentSelectedWaypoint != nil {
		currentSelectedWaypointInfo := currentSelectedWaypoint.(waypointInfoItem)
		m.WaypointInfoList.AdditionalShortHelpKeys = initShortWaypointInfoListKeyMap(NoPopUp, currentSelectedWaypointInfo.WaypointIsSealed)
	}

	return nil
}

// ----------------------------------
//
//	returns an AdditionalShortHelpKeys func whose bindings are determined by
//	popUpType and, when popUpType is NoPopUp, by isSealed; for NoPopUp it
//	shows ? (help) plus the primary action key relevant to the item's sealed
//	state, keeping the status bar compact; for RebindPopUp and ReforgePopUp
//	it shows only the popup-relevant keys: enter (submit), esc (close), and
//	q/ctrl+c (quit)
//
// ----------------------------------
func initShortWaypointInfoListKeyMap(popUpType string, isSealed bool) func() []key.Binding {
	switch popUpType {
	case NoPopUp:
		if isSealed {
			return func() []key.Binding {
				return []key.Binding{
					key.NewBinding(key.WithKeys("?"), key.WithHelp("?", i18n.LANGUAGEMAPPING.WaypointUIHelpKeyHelp)),
					key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
				}
			}
		} else {
			return func() []key.Binding {
				return []key.Binding{
					key.NewBinding(key.WithKeys("?"), key.WithHelp("?", i18n.LANGUAGEMAPPING.WaypointUIHelpKeyHelp)),
					key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.WaypointNavigateKeyHelp)),
					key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
				}
			}
		}
	case RebindPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.WaypointRebindKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.WaypointClosePopUp)),
				key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
			}
		}
	case ReforgePopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.WaypointReforgeKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.WaypointClosePopUp)),
				key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp)),
			}
		}
	case HelpPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.WaypointClosePopUp)),
			}
		}
	}

	return func() []key.Binding { return []key.Binding{} }
}

// ----------------------------------
//
//	builds the full key binding reference and stores it as rendered text in
//	the help viewport; each entry is formatted as keybinding → action →
//	description, all coloured with the purple palette; called once during
//	model initialisation so the content is ready before the first render
//
// ----------------------------------
func initWaypointHelpViewport(m *WaypointInteractiveModel) {
	helpList := []HelpListItem{
		{Keybinding: "?", Action: i18n.LANGUAGEMAPPING.WaypointUIHelpKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUIHelpKeyHelpDescription},
		{Keybinding: "↑/k", Action: i18n.LANGUAGEMAPPING.WaypointUIUpKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUIUpKeyHelpDescription},
		{Keybinding: "↓/j", Action: i18n.LANGUAGEMAPPING.WaypointUIDownKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUIDownKeyHelpDescription},
		{Keybinding: "enter", Action: i18n.LANGUAGEMAPPING.WaypointNavigateKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointNavigateKeyHelpDescription},
		{Keybinding: "u/U", Action: i18n.LANGUAGEMAPPING.WaypointUnsealKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUnsealKeyHelpDescription},
		{Keybinding: "r", Action: i18n.LANGUAGEMAPPING.WaypointRebindKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointRebindKeyHelpDescription},
		{Keybinding: "R", Action: i18n.LANGUAGEMAPPING.WaypointReforgeKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointReforgeKeyHelpDescription},
		{Keybinding: "y/Y", Action: i18n.LANGUAGEMAPPING.WaypointNameCopyPathCopyKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointNameCopyPathCopyKeyHelpDescription},
		{Keybinding: "backspace", Action: i18n.LANGUAGEMAPPING.WaypointDestroyKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointDestroyKeyHelpDescription},
		{Keybinding: "ctrl+y", Action: i18n.LANGUAGEMAPPING.WaypointCopyFromInputValueKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointCopyFromInputValueKeyHelpDescription},
		{Keybinding: "ctrl+p", Action: i18n.LANGUAGEMAPPING.WaypointPasteIntoInputValueKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointPasteIntoInputValueKeyHelpDescription},
		{Keybinding: "esc", Action: i18n.LANGUAGEMAPPING.WaypointClosePopUp, Description: i18n.LANGUAGEMAPPING.WaypointClosePopUpDescription},
		{Keybinding: "q/esc/ctrl+c", Action: i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelp, Description: i18n.LANGUAGEMAPPING.WaypointUIQuitKeyHelpDescription},
	}

	var content strings.Builder
	for _, helpitem := range helpList {
		keybinding := style.RenderStringWithColor(fmt.Sprintf("[%s]", helpitem.Keybinding), style.ColorPurpleVibrant, false)
		action := style.RenderStringWithColor(helpitem.Action, style.ColorPurpleSoft, false)
		description := style.RenderStringWithColor(helpitem.Description, style.ColorPurpleSoft, true)
		content.WriteString(keybinding)
		content.WriteRune('\n')
		content.WriteString(action)
		content.WriteRune('\n')
		content.WriteString(description)
		content.WriteRune('\n')
		content.WriteRune('\n')
	}

	m.WaypointHelpViewport.SetContent(content.String())
}
