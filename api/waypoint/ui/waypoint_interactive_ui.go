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
	"go.etcd.io/bbolt"
)

// ----------------------------------
//
//	bubbletea model that drives the interactive waypoint selection UI;
//	holds all display and navigation state across Init / Update / View
//
// ----------------------------------
type WaypointInteractiveModel struct {
	SelectedWaypointPath           string
	SelectedWaypointName           string
	IsQuit                         bool
	ShowPopUp                      atomic.Bool
	PopUpType                      string
	IsTypingMode                   atomic.Bool
	ErrMessage                     error
	WaypointInfoList               list.Model
	WaypointInfoListCursorPosition int
	WaypointHelpViewport           viewport.Model
	WaypointPopUpModel             interface{}
	Width                          int
	Height                         int
	BboltReadDb                    *bbolt.DB
	IsRenderInit                   atomic.Bool
}

// ----------------------------------
//
//	allocates a WaypointInteractiveModel with safe zero-value defaults;
//	list setup is deferred to the first WindowSizeMsg so that correct
//	terminal dimensions are available when building item layouts
//
// ----------------------------------
func initWaypointInteractiveModel(bboltReadDb *bbolt.DB) *WaypointInteractiveModel {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = false
	vp.SetHorizontalStep(1)

	waypointInteractiveModel := WaypointInteractiveModel{
		SelectedWaypointPath:           "",
		IsQuit:                         false,
		PopUpType:                      NoPopUp,
		WaypointInfoListCursorPosition: 0,
		WaypointHelpViewport:           vp,
		BboltReadDb:                    bboltReadDb,
	}
	waypointInteractiveModel.ShowPopUp.Store(false)
	waypointInteractiveModel.IsTypingMode.Store(false)
	initWaypointHelpViewport(&waypointInteractiveModel)

	return &waypointInteractiveModel
}

// ----------------------------------
//
//	starts the bubbletea program and blocks until the user selects a
//	waypoint or quits; returns the selected waypoint path, name, and
//	any error encountered during the program run or after selection;
//	note: live terminal window resizing is not supported — rift is a
//	CLI-first tool and the interactive UI is a quality-of-life convenience,
//	not a full TUI application
//
// ----------------------------------
func RunWaypointInteractive(bboltReadDb *bbolt.DB) (string, string, error) {
	waypointInteractiveModel := initWaypointInteractiveModel(bboltReadDb)
	// route program output to stderr so stdout stays clean for callers
	p := tea.NewProgram(waypointInteractiveModel, tea.WithOutput(os.Stderr))
	result, err := p.Run()
	if err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointInteractiveError, err), style.ColorError, false)
		return "", "", fmt.Errorf("%s", errorMessage)
	}

	// safe: Run always returns the model passed to NewProgram
	final := result.(*WaypointInteractiveModel)
	if final.ErrMessage != nil {
		return "", "", final.ErrMessage
	}

	// return selected path
	if final.SelectedWaypointPath == "" || final.SelectedWaypointName == "" {
		return "", "", nil
	}
	return final.SelectedWaypointPath, final.SelectedWaypointName, nil
}

// ----------------------------------
//
//	satisfies the tea.Model interface; no startup commands are needed
//
// ----------------------------------
func (m *WaypointInteractiveModel) Init() tea.Cmd {
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
//	sealed state of the selected waypoint
//
// ----------------------------------
func (m *WaypointInteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			initWaypointInfoListModel(m)
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
				// refresh the help key map to reflect the sealed state of the new selection
				if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
					m.WaypointInfoList.AdditionalShortHelpKeys = initShortWaypointInfoListKeyMap(m.PopUpType, i.WaypointIsSealed)
				}
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

		// refresh the help key map to reflect the sealed state of the new selection
		if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
			m.WaypointInfoList.AdditionalShortHelpKeys = initShortWaypointInfoListKeyMap(m.PopUpType, i.WaypointIsSealed)
		}

		return m, tea.Batch(cmds...)

	}
	return m, tea.Batch(cmds...)
}

// ----------------------------------
//
//	renders the waypoint list; returns an empty view once the user
//	has quit to avoid a final flash of stale content
//
// ----------------------------------
func (m *WaypointInteractiveModel) View() tea.View {
	if m.IsQuit {
		return tea.NewView("")
	}
	return tea.NewView(renderWaypointInteractiveUIView(m))
}
