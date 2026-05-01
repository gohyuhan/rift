package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Reads every entry in the ritual bucket and returns a slice of ritualInfo
//	records. Uses a read-only View transaction; the read connection is closed
//	immediately after View completes so it does not block any concurrent write
//	connection. Corrupted proto entries are collected during the View and
//	recorded via RecordCorruptedRitualInfo in a follow-up Update transaction.
//
// ----------------------------------
func getAllRitualInfo() ([]ritualInfo, error) {
	var ritualsInfo []ritualInfo
	var corruptedRitualName []string
	ritualCorrupted := false

	// open DB so we can read ritual records
	viewErr := func() error {
		bboltReadDb, bboltReadDbErr := db.OpenReadDB()
		if bboltReadDbErr != nil {
			return bboltReadDbErr
		}
		defer db.CloseDB(bboltReadDb)

		return bboltReadDb.View(func(tx *bbolt.Tx) error {
			// ensure the ritual bucket exists before iterating
			ritualBucket := tx.Bucket(db.RitualBucket)
			if ritualBucket == nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RitualBucketNotFoundError, style.ColorError, false))
			}

			// walk every key-value pair in the bucket
			retrieveError := ritualBucket.ForEach(func(k, v []byte) error {
				// deserialize the stored proto; capture the name, set the flag, and
				// return a sentinel error to stop ForEach — recording is deferred to
				// a separate Update after the View transaction completes
				existingRitual := &pb.Ritual{}
				protoErr := proto.Unmarshal(v, existingRitual)

				// skip corrupted data
				if protoErr != nil {
					ritualCorrupted = true
					corruptedRitualName = append(corruptedRitualName, string(k))
					return nil
				}

				// construct the ritual info type
				info := ritualInfo{
					RitualName: string(k),
					RitualDesc: existingRitual.RitualDesc,
				}

				ritualsInfo = append(ritualsInfo, info)

				return nil
			})

			// wrap ForEach failure into a user-facing message
			if retrieveError != nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftRitualRetrieveAllError, style.ColorError, false))
			}

			return nil
		})
	}()

	if ritualCorrupted {
		viewErr = apiUtils.RecordCorruptedRitualInfo(corruptedRitualName)
	}

	return ritualsInfo, viewErr
}

// ----------------------------------
//
//	builds or rebuilds the bubbles list model on the ScrollInteractiveModel;
//	preserves the previously selected ritual by name when re-initialising,
//	and falls back to the stored cursor position when no match is found;
//	called after the first WindowSizeMsg (so valid dimensions are available)
//	and again whenever the ritual list changes structurally (e.g. after a forget)
//
// ----------------------------------
func initRitualInfoListModel(m *ScrollInteractiveModel) error {
	previousSelectedRitual := m.RitualInfoList.SelectedItem()
	selectedRitualCursorPosition := -1

	latestRitualInfoArray := []list.Item{}

	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - ListTitleHorizontalPadding

	allRitualsInfo, err := getAllRitualInfo()
	if err != nil {
		return err
	}

	if previousSelectedRitual != nil {
		// a previous selection exists; scan the fresh list to find a name match
		// so the cursor lands on the same ritual after re-initialisation
		previousSelectedRitualInfo := previousSelectedRitual.(ritualInfoItem)
		for index, ritual := range allRitualsInfo {
			if ritual.RitualName == previousSelectedRitualInfo.RitualName {
				selectedRitualCursorPosition = index
			}
			latestRitualInfoArray = append(latestRitualInfoArray, ritualInfoItem(ritual))
		}
	} else {
		// no previous selection; convert all records to list items
		for _, ritual := range allRitualsInfo {
			latestRitualInfoArray = append(latestRitualInfoArray, ritualInfoItem(ritual))
		}
	}

	m.RitualInfoList = list.New(latestRitualInfoArray, ritualInfoDelegate{}, m.Width, m.Height)
	m.RitualInfoList.SetShowPagination(false)
	m.RitualInfoList.SetShowStatusBar(false)
	m.RitualInfoList.SetFilteringEnabled(false)
	m.RitualInfoList.SetShowFilter(false)
	m.RitualInfoList.SetShowHelp(true)

	// truncate the title to prevent overflow when the terminal is narrow
	m.RitualInfoList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.RitualInfoListTitle, titleWidthLimit, "...")
	m.RitualInfoList.Styles.Title = style.NewStyle.Bold(true)
	m.RitualInfoList.Styles.PaginationStyle = style.NewStyle
	m.RitualInfoList.Styles.TitleBar = style.NewStyle
	m.RitualInfoList.Styles.HelpStyle = style.NewStyle
	m.RitualInfoList.Help.Styles.ShortKey = style.NewStyle.Foreground(style.ColorBlueMuted).Bold(true)
	m.RitualInfoList.Help.Styles.ShortDesc = style.NewStyle.Foreground(style.ColorBlueMuted)
	m.RitualInfoList.Help.Styles.ShortSeparator = style.NewStyle.Foreground(style.ColorBlueGrayMuted)
	// clear the default key map so only our custom bindings appear in the help bar
	m.RitualInfoList.KeyMap = list.KeyMap{}
	m.RitualInfoList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.RitualUIQuitKeyHelp)),
		}
	}

	if selectedRitualCursorPosition >= 0 {
		// restore the cursor to the previously selected ritual by name match
		m.RitualInfoList.Select(selectedRitualCursorPosition)
		m.RitualInfoListCursorPosition = selectedRitualCursorPosition
	} else {
		// clamp the stored position to the last item if the list has shrunk
		if m.RitualInfoListCursorPosition > len(m.RitualInfoList.Items())-1 {
			m.RitualInfoList.Select(len(m.RitualInfoList.Items()) - 1)
			m.RitualInfoListCursorPosition = len(m.RitualInfoList.Items()) - 1
		} else {
			m.RitualInfoList.Select(m.RitualInfoListCursorPosition)
		}
	}

	// update the help key map
	currentSelectedRitual := m.RitualInfoList.SelectedItem()
	if currentSelectedRitual != nil {
		m.RitualInfoList.AdditionalShortHelpKeys = initShortRitualInfoListKeyMap(NoPopUp)
	}

	return nil
}

// ----------------------------------
//
//	returns an AdditionalShortHelpKeys func whose bindings are determined by
//	popUpType; for NoPopUp it shows ? (help), enter (cast), and q/esc/ctrl+c
//	(quit), keeping the status bar compact; unrecognised popup types return
//	an empty binding list
//
// ----------------------------------
func initShortRitualInfoListKeyMap(popUpType string) func() []key.Binding {
	switch popUpType {
	case NoPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("?"), key.WithHelp("?", i18n.LANGUAGEMAPPING.RitualUIHelpKeyHelp)),
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.RitualInvokeKeyHelp)),
				key.NewBinding(key.WithKeys("n"), key.WithHelp("n", i18n.LANGUAGEMAPPING.RitualInscribeKeyHelp)),
				key.NewBinding(key.WithKeys("e"), key.WithHelp("e", i18n.LANGUAGEMAPPING.RitualReinscribeKeyHelp)),
				key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.RitualUIQuitKeyHelp)),
			}
		}
	case HelpPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.RitualClosePopUp)),
			}
		}
	case InscribePopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.RitualUIInscribeKeyHelp)),
				key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", i18n.LANGUAGEMAPPING.RitualUINextInputKeyHelp)),
				key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", i18n.LANGUAGEMAPPING.RitualUIPreviousInputKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.RitualClosePopUp)),
			}
		}
	case InvokeLocationOptionPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.RitualUIChooseInvokeLocationKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.RitualClosePopUp)),
			}
		}
	case InvokeWaypointLocationOptionPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.RitualUIChooseWaypointInvokeLocationKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.RitualClosePopUp)),
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
func initRitualHelpViewport(m *ScrollInteractiveModel) {
	helpList := []HelpListItem{
		{Keybinding: "?", Action: i18n.LANGUAGEMAPPING.RitualUIHelpKeyHelp, Description: i18n.LANGUAGEMAPPING.RitualUIHelpKeyHelpDescription},
		{Keybinding: "↑/k", Action: i18n.LANGUAGEMAPPING.RitualUIUpKeyHelp, Description: i18n.LANGUAGEMAPPING.RitualUIUpKeyHelpDescription},
		{Keybinding: "↓/j", Action: i18n.LANGUAGEMAPPING.RitualUIDownKeyHelp, Description: i18n.LANGUAGEMAPPING.RitualUIDownKeyHelpDescription},
		{Keybinding: "n/N", Action: i18n.LANGUAGEMAPPING.RitualInscribeKeyHelp, Description: i18n.LANGUAGEMAPPING.RitualInscribeKeyHelpDescription},
		{Keybinding: "e/E", Action: i18n.LANGUAGEMAPPING.RitualReinscribeKeyHelp, Description: i18n.LANGUAGEMAPPING.RitualReinscribeKeyHelpDescription},
		{Keybinding: "enter", Action: i18n.LANGUAGEMAPPING.RitualInvokeKeyHelp, Description: i18n.LANGUAGEMAPPING.RitualInvokeKeyHelpDescription},
		{Keybinding: "backspace", Action: i18n.LANGUAGEMAPPING.RitualForgetKey, Description: i18n.LANGUAGEMAPPING.RitualForgetKeyDescription},
		{Keybinding: "esc", Action: i18n.LANGUAGEMAPPING.RitualClosePopUp, Description: i18n.LANGUAGEMAPPING.RitualClosePopUpDescription},
		{Keybinding: "q/esc/ctrl+c", Action: i18n.LANGUAGEMAPPING.RitualUIQuitKeyHelp, Description: i18n.LANGUAGEMAPPING.RitualUIQuitKeyHelpDescription},
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

	m.RitualHelpViewport.SetContent(content.String())
}
