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
//	reads every entry in the spell bucket and returns a slice of
//	spellInfo records; uses a read-only View transaction so any
//	corruption writes are deferred to a separate Update after View
//	completes; corrupted proto entries are collected and recorded via
//	RecordCorruptedSpellInfo before the error is returned to the caller
//
// ----------------------------------
func getAllSpellInfo(bboltReadDb *bbolt.DB) ([]spellInfo, error) {
	var spellsInfo []spellInfo
	var corruptedSpellName []string
	spellCorrupted := false

	viewErr := bboltReadDb.View(func(tx *bbolt.Tx) error {
		// ensure the spell bucket exists before iterating
		spellBucket := tx.Bucket(db.SpellBucket)
		if spellBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellBucketNotFoundError, style.ColorError, false))
		}

		// walk every key-value pair in the bucket
		retrieveError := spellBucket.ForEach(func(k, v []byte) error {
			// deserialize the stored proto; capture the name, set the flag, and
			// return a sentinel error to stop ForEach — recording is deferred to
			// a separate Update after the View transaction completes
			existingSpell := &pb.Spell{}
			protoErr := proto.Unmarshal(v, existingSpell)

			// skip corrupted data
			if protoErr != nil {
				spellCorrupted = true
				corruptedSpellName = append(corruptedSpellName, string(k))
				return nil
			}

			// construct the spell info type
			info := spellInfo{
				SpellName:    string(k),
				SpellCommand: existingSpell.SpellCommand,
			}

			spellsInfo = append(spellsInfo, info)

			return nil
		})

		// wrap ForEach failure into a user-facing message
		if retrieveError != nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RiftSpellRetrieveAllError, style.ColorError, false))
		}

		return nil
	})

	if spellCorrupted {
		viewErr = apiUtils.RecordCorruptedSpellInfo(corruptedSpellName)
	}

	return spellsInfo, viewErr
}

// ----------------------------------
//
//	builds or rebuilds the bubbles list model on the SpellbookInteractiveModel;
//	preserves the previously selected spell by name when re-initialising,
//	and falls back to the stored cursor position when no match is found;
//	called after the first WindowSizeMsg (so valid dimensions are available)
//	and again whenever the spell list changes structurally (e.g. after a forget)
//
// ----------------------------------
func initSpellInfoListModel(m *SpellbookInteractiveModel) error {
	previousSelectedSpell := m.SpellInfoList.SelectedItem()
	selectedSpellCursorPosition := -1

	latestSpellInfoArray := []list.Item{}

	titleWidthLimit := m.Width - ListItemOrTitleWidthPad - ListTitleHorizontalPadding

	allSpellsInfo, err := getAllSpellInfo(m.BboltReadDb)
	if err != nil {
		return err
	}

	if previousSelectedSpell != nil {
		// a previous selection exists; scan the fresh list to find a name match
		// so the cursor lands on the same spell after re-initialisation
		previousSelectedSpellInfo := previousSelectedSpell.(spellInfoItem)
		for index, spell := range allSpellsInfo {
			if spell.SpellName == previousSelectedSpellInfo.SpellName {
				selectedSpellCursorPosition = index
			}
			latestSpellInfoArray = append(latestSpellInfoArray, spellInfoItem(spell))
		}
	} else {
		// no previous selection; convert all records to list items
		for _, spell := range allSpellsInfo {
			latestSpellInfoArray = append(latestSpellInfoArray, spellInfoItem(spell))
		}
	}

	m.SpellInfoList = list.New(latestSpellInfoArray, spellInfoDelegate{}, m.Width, m.Height)
	m.SpellInfoList.SetShowPagination(false)
	m.SpellInfoList.SetShowStatusBar(false)
	m.SpellInfoList.SetFilteringEnabled(false)
	m.SpellInfoList.SetShowFilter(false)
	m.SpellInfoList.SetShowHelp(true)

	// truncate the title to prevent overflow when the terminal is narrow
	m.SpellInfoList.Title = ansi.Truncate(i18n.LANGUAGEMAPPING.SpellInfoListTitle, titleWidthLimit, "...")
	m.SpellInfoList.Styles.Title = style.NewStyle.Bold(true)
	m.SpellInfoList.Styles.PaginationStyle = style.NewStyle
	m.SpellInfoList.Styles.TitleBar = style.NewStyle
	m.SpellInfoList.Styles.HelpStyle = style.NewStyle
	m.SpellInfoList.Help.Styles.ShortKey = style.NewStyle.Foreground(style.ColorBlueMuted).Bold(true)
	m.SpellInfoList.Help.Styles.ShortDesc = style.NewStyle.Foreground(style.ColorBlueMuted)
	m.SpellInfoList.Help.Styles.ShortSeparator = style.NewStyle.Foreground(style.ColorBlueGrayMuted)
	// clear the default key map so only our custom bindings appear in the help bar
	m.SpellInfoList.KeyMap = list.KeyMap{}
	m.SpellInfoList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.SpellUIQuitKeyHelp)),
		}
	}

	if selectedSpellCursorPosition >= 0 {
		// restore the cursor to the previously selected spell by name match
		m.SpellInfoList.Select(selectedSpellCursorPosition)
		m.SpellInfoListCursorPosition = selectedSpellCursorPosition
	} else {
		// clamp the stored position to the last item if the list has shrunk
		if m.SpellInfoListCursorPosition > len(m.SpellInfoList.Items())-1 {
			m.SpellInfoList.Select(len(m.SpellInfoList.Items()) - 1)
			m.SpellInfoListCursorPosition = len(m.SpellInfoList.Items()) - 1
		} else {
			m.SpellInfoList.Select(m.SpellInfoListCursorPosition)
		}
	}

	// update the help key map
	currentSelectedSpell := m.SpellInfoList.SelectedItem()
	if currentSelectedSpell != nil {
		m.SpellInfoList.AdditionalShortHelpKeys = initShortSpellInfoListKeyMap(NoPopUp)
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
func initShortSpellInfoListKeyMap(popUpType string) func() []key.Binding {
	switch popUpType {
	case NoPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("?"), key.WithHelp("?", i18n.LANGUAGEMAPPING.SpellUIHelpKeyHelp)),
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.SpellCastKeyHelp)),
				key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc/ctrl+c", i18n.LANGUAGEMAPPING.SpellUIQuitKeyHelp)),
			}
		}
	case HelpPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.SpellClosePopUp)),
			}
		}
	case LearnPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.SpellUILearnKeyHelp)),
				key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", i18n.LANGUAGEMAPPING.SpellUINextInputKeyHelp)),
				key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", i18n.LANGUAGEMAPPING.SpellUIPreviousInputKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.SpellClosePopUp)),
			}
		}
	case CastLocationOptionPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.SpellUIChooseCastLocationKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.SpellClosePopUp)),
			}
		}
	case CastWaypointLocationOptionPopUp:
		return func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", i18n.LANGUAGEMAPPING.SpellUIChooseWaypointCastLocationKeyHelp)),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", i18n.LANGUAGEMAPPING.SpellClosePopUp)),
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
func initSpellHelpViewport(m *SpellbookInteractiveModel) {
	helpList := []HelpListItem{
		{Keybinding: "?", Action: i18n.LANGUAGEMAPPING.SpellUIHelpKeyHelp, Description: i18n.LANGUAGEMAPPING.SpellUIHelpKeyHelpDescription},
		{Keybinding: "↑/k", Action: i18n.LANGUAGEMAPPING.SpellUIUpKeyHelp, Description: i18n.LANGUAGEMAPPING.SpellUIUpKeyHelpDescription},
		{Keybinding: "↓/j", Action: i18n.LANGUAGEMAPPING.SpellUIDownKeyHelp, Description: i18n.LANGUAGEMAPPING.SpellUIDownKeyHelpDescription},
		{Keybinding: "n/N", Action: i18n.LANGUAGEMAPPING.SpellLearnKeyHelp, Description: i18n.LANGUAGEMAPPING.SpellLearnKeyHelpDescription},
		{Keybinding: "enter", Action: i18n.LANGUAGEMAPPING.SpellCastKeyHelp, Description: i18n.LANGUAGEMAPPING.SpellCastKeyHelpDescription},
		{Keybinding: "backspace", Action: i18n.LANGUAGEMAPPING.SpellForgetKey, Description: i18n.LANGUAGEMAPPING.SpellForgetKeyDescription},
		{Keybinding: "esc", Action: i18n.LANGUAGEMAPPING.SpellClosePopUp, Description: i18n.LANGUAGEMAPPING.SpellClosePopUpDescription},
		{Keybinding: "q/esc/ctrl+c", Action: i18n.LANGUAGEMAPPING.SpellUIQuitKeyHelp, Description: i18n.LANGUAGEMAPPING.SpellUIQuitKeyHelpDescription},
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

	m.SpellHelpViewport.SetContent(content.String())
}
