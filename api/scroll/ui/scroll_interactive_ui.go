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
)

// ----------------------------------
//
//	bubbletea model that drives the interactive scroll selection UI;
//	holds all display and navigation state across Init / Update / View
//
// ----------------------------------
type ScrollInteractiveModel struct {
	SelectedRitualName           string
	RitualInvokePath             string
	IsQuit                       bool
	ShowPopUp                    atomic.Bool
	PopUpType                    string
	IsTypingMode                 atomic.Bool
	ErrMessage                   error
	RitualInfoList               list.Model
	RitualInfoListCursorPosition int
	RitualHelpViewport           viewport.Model
	RitualPopUpModel             interface{}
	Width                        int
	Height                       int
	IsRenderInit                 atomic.Bool
}

// ----------------------------------
//
//	raw data record for a single ritual as read from the database;
//	mirrors the proto fields relevant to the list UI
//
// ----------------------------------
type ritualInfo struct {
	RitualName string
	RitualDesc string
}

// ----------------------------------
//
//	allocates a ScrollInteractiveModel with safe zero-value defaults;
//	list setup is deferred to the first WindowSizeMsg so that correct
//	terminal dimensions are available when building item layouts
//
// ----------------------------------
func initScrollInteractiveModel() *ScrollInteractiveModel {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = false
	vp.SetHorizontalStep(1)

	scrollInteractiveModel := ScrollInteractiveModel{
		SelectedRitualName:           "",
		RitualInvokePath:             "",
		IsQuit:                       false,
		PopUpType:                    NoPopUp,
		RitualInfoListCursorPosition: 0,
		RitualHelpViewport:           vp,
	}
	scrollInteractiveModel.ShowPopUp.Store(false)
	scrollInteractiveModel.IsTypingMode.Store(false)
	initRitualHelpViewport(&scrollInteractiveModel)

	return &scrollInteractiveModel
}

// ----------------------------------
//
//	starts the bubbletea program and blocks until the user selects a
//	ritual or quits; returns the selected ritual command, name, and any
//	error encountered during the program run or after selection;
//	note: live terminal window resizing is not supported — rift is a
//	CLI-first tool and the interactive UI is a quality-of-life convenience,
//	not a full TUI application
//
// ----------------------------------
func RunScrollInteractive() (string, string, error) {
	scrollInteractiveModel := initScrollInteractiveModel()
	// route program output to stderr so stdout stays clean for callers
	p := tea.NewProgram(scrollInteractiveModel, tea.WithOutput(os.Stderr))
	result, err := p.Run()
	if err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.ScrollInteractiveError, err), style.ColorError, false)
		return "", "", fmt.Errorf("%s", errorMessage)
	}

	// safe: Run always returns the model passed to NewProgram
	final := result.(*ScrollInteractiveModel)
	if final.ErrMessage != nil {
		return "", "", final.ErrMessage
	}

	// if user quit without selecting a ritual, return empty values
	if final.SelectedRitualName == "" {
		return "", "", nil
	}

	// if user selected a ritual but the cast path is empty, default to the current working directory
	if final.RitualInvokePath == "" && final.SelectedRitualName != "" {
		path, err := utils.GetCWD()
		if err != nil {
			return "", "", err
		}
		final.RitualInvokePath = path
	}

	return final.SelectedRitualName, final.RitualInvokePath, nil
}

// ----------------------------------
//
//	satisfies the tea.Model interface; no startup commands are needed
//
// ----------------------------------
func (m *ScrollInteractiveModel) Init() tea.Cmd {
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
//	the selected ritual
//
// ----------------------------------
func (m *ScrollInteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			initRitualInfoListModel(m)
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
				m.RitualInfoList.AdditionalShortHelpKeys = initShortRitualInfoListKeyMap(m.PopUpType)
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
		m.RitualInfoList.AdditionalShortHelpKeys = initShortRitualInfoListKeyMap(m.PopUpType)

		return m, tea.Batch(cmds...)

	}
	return m, tea.Batch(cmds...)
}

// ----------------------------------
//
//	renders the scroll list; returns an empty view once the user
//	has quit to avoid a final flash of stale content
//
// ----------------------------------
func (m *ScrollInteractiveModel) View() tea.View {
	if m.IsQuit {
		return tea.NewView("")
	}
	return tea.NewView(renderScrollInteractiveUIView(m))
}
