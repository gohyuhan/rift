package ui

import (
	"fmt"
	"os"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"go.etcd.io/bbolt"
)

// ----------------------------------
//
//	bubbletea model that drives the interactive spellbook selection UI;
//	holds all display and navigation state across Init / Update / View
//
// ----------------------------------
type SpellbookInteractiveModel struct {
	SelectedSpellName           string
	SpellCastPath               string
	IsQuit                      bool
	ShowPopUp                   atomic.Bool
	PopUpType                   string
	IsTypingMode                atomic.Bool
	ErrMessage                  error
	SpellInfoList               list.Model
	SpellInfoListCursorPosition int
	SpellHelpViewport           viewport.Model
	SpellPopUpModel             interface{}
	Width                       int
	Height                      int
	BboltDb                     *bbolt.DB
	IsRenderInit                atomic.Bool
}

// ----------------------------------
//
//	raw data record for a single spell as read from the database;
//	mirrors the proto fields relevant to the list UI
//
// ----------------------------------
type spellInfo struct {
	SpellName    string
	SpellCommand []string
}

// ----------------------------------
//
//	allocates a SpellbookInteractiveModel with safe zero-value defaults;
//	list setup is deferred to the first WindowSizeMsg so that correct
//	terminal dimensions are available when building item layouts
//
// ----------------------------------
func initSpellbookInteractiveModel(bboltDb *bbolt.DB) *SpellbookInteractiveModel {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = false
	vp.SetHorizontalStep(1)

	spellbookInteractiveModel := SpellbookInteractiveModel{
		SelectedSpellName:           "",
		SpellCastPath:               "",
		IsQuit:                      false,
		PopUpType:                   NoPopUp,
		SpellInfoListCursorPosition: 0,
		SpellHelpViewport:           vp,
		BboltDb:                     bboltDb,
	}
	spellbookInteractiveModel.ShowPopUp.Store(false)
	spellbookInteractiveModel.IsTypingMode.Store(false)
	initSpellHelpViewport(&spellbookInteractiveModel)

	return &spellbookInteractiveModel
}

// ----------------------------------
//
//	starts the bubbletea program and blocks until the user selects a
//	spell or quits; returns the selected spell command, name, and any
//	error encountered during the program run or after selection;
//	note: live terminal window resizing is not supported — rift is a
//	CLI-first tool and the interactive UI is a quality-of-life convenience,
//	not a full TUI application
//
// ----------------------------------
func RunSpellbookInteractive(bboltDb *bbolt.DB) (string, string, error) {
	spellbookInteractiveModel := initSpellbookInteractiveModel(bboltDb)
	// route program output to stderr so stdout stays clean for callers
	p := tea.NewProgram(spellbookInteractiveModel, tea.WithOutput(os.Stderr))
	result, err := p.Run()
	if err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.SpellbookInteractiveError, err), style.ColorError, false)
		return "", "", fmt.Errorf("%s", errorMessage)
	}

	// safe: Run always returns the model passed to NewProgram
	final := result.(*SpellbookInteractiveModel)
	if final.ErrMessage != nil {
		return "", "", final.ErrMessage
	}

	// if user quit without selecting a spell, return empty values
	if final.SelectedSpellName == "" {
		return "", "", nil
	}

	// if user selected a spell but the cast path is empty, default to the current working directory;z
	if final.SpellCastPath == "" && final.SelectedSpellName != "" {
		path, err := utils.GetCWD()
		if err != nil {
			return "", "", err
		}
		final.SpellCastPath = path
	}

	return final.SelectedSpellName, final.SpellCastPath, nil
}

// ----------------------------------
//
//	satisfies the tea.Model interface; no startup commands are needed
//
// ----------------------------------
func (m *SpellbookInteractiveModel) Init() tea.Cmd {
	return nil
}

// ----------------------------------
//
//	handles terminal resize and key events; list initialisation is deferred
//	until the first WindowSizeMsg so that valid dimensions are available for
//	layout calculation; ctrl+c/q always quit; esc closes the active popup
//	(if any) or quits when no popup is open; all other key events are
//	dispatched to handleTypingInteraction when a text input popup is active,
//	or to handleNonTypingInteraction otherwise; the help key map is refreshed
//	after every key event to stay consistent with the current popup type and
//	the selected spell
//
// ----------------------------------
func (m *SpellbookInteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// initialize list components once, immediately after the first window resize;
		// valid dimensions are required to calculate item layouts (specifically text
		// truncation) — initializing earlier would cause the UI layout to break
		if m.IsRenderInit.CompareAndSwap(false, true) {
			initSpellInfoListModel(m)
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.IsQuit = true
			return m, tea.Quit
		case "esc":
			if m.ShowPopUp.Load() {
				m.ShowPopUp.Store(false)
				m.IsTypingMode.Store(false)
				m.PopUpType = NoPopUp
				m.SpellInfoList.AdditionalShortHelpKeys = initShortSpellInfoListKeyMap(m.PopUpType)
				return m, nil
			}
			m.IsQuit = true
			return m, tea.Quit
		}
		if m.IsTypingMode.Load() {
			m, cmd = handleTypingInteraction(m, msg)
		} else {
			m, cmd = handleNonTypingInteraction(m, msg)
		}

		cmds = append(cmds, cmd)
		m.SpellInfoList.AdditionalShortHelpKeys = initShortSpellInfoListKeyMap(m.PopUpType)

		return m, tea.Batch(cmds...)

	}
	return m, tea.Batch(cmds...)
}

// ----------------------------------
//
//	renders the spellbook list; returns an empty view once the user
//	has quit to avoid a final flash of stale content
//
// ----------------------------------
func (m *SpellbookInteractiveModel) View() tea.View {
	if m.IsQuit {
		return tea.NewView("")
	}
	return tea.NewView(renderSpellbookInteractiveUIView(m))
}
