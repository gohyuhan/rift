package ui

import (
	"fmt"
	"os"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/api/waypoint/features"
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
	ErrMessage                     error
	WaypointInfoList               list.Model
	WaypointInfoListCursorPosition int
	Width                          int
	Height                         int
	BboltDb                        *bbolt.DB
	IsRenderInit                   atomic.Bool
}

// ----------------------------------
//
//	raw data record for a single waypoint as read from the database;
//	mirrors the proto fields relevant to the list UI
//
// ----------------------------------
type waypointInfo struct {
	WaypointName         string
	WaypointPath         string
	WaypointIsSealed     bool
	WaypointSealedReason string
}

// ----------------------------------
//
//	allocates a WaypointInteractiveModel with safe zero-value defaults;
//	list setup is deferred to the first WindowSizeMsg so that correct
//	terminal dimensions are available when building item layouts
//
// ----------------------------------
func initWaypointInteractiveModel(bboltDb *bbolt.DB) *WaypointInteractiveModel {
	waypointInteractiveModel := WaypointInteractiveModel{
		SelectedWaypointPath:           "",
		IsQuit:                         false,
		WaypointInfoListCursorPosition: 0,
		BboltDb:                        bboltDb,
	}

	return &waypointInteractiveModel
}

// ----------------------------------
//
//	starts the bubbletea program and blocks until the user selects a
//	waypoint or quits; returns the selected waypoint path, name, and
//	any error encountered during the program run or after selection
//
// ----------------------------------
func RunWaypointInteractive(bboltDb *bbolt.DB) (string, string, error) {
	waypointInteractiveModel := initWaypointInteractiveModel(bboltDb)
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
//	handles terminal resize, cursor navigation, selection, quit, and unseal;
//	list component initialisation is deferred until the first resize
//	event so that valid dimensions are available for layout calculation;
//	enter navigates to the selected waypoint (no-op if sealed);
//	u/U attempts to unseal the selected waypoint and rebuilds the list
//	(if the path is still missing the waypoint will be re-sealed immediately);
//	backspace permanently destroys the selected waypoint and rebuilds the list;
//	j/k and ↓/↑ move the cursor and refresh the help key map to reflect
//	the sealed state of the newly focused item
//
// ----------------------------------
func (m *WaypointInteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter":
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				// sealed waypoints are non-navigable; silently ignore enter
				if !i.WaypointIsSealed {
					m.SelectedWaypointPath = i.WaypointPath
					m.SelectedWaypointName = i.WaypointName
					m.IsQuit = true
					return m, tea.Quit
				}
				return m, nil
			}
		case "backspace":
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				// perform a destroy on the selected waypoint
				destoryErr := features.DestroyDiscoveredWaypoint(m.BboltDb, i.WaypointName, false)
				if destoryErr != nil {
					m.ErrMessage = destoryErr
					m.IsQuit = true
					return m, tea.Quit
				}
				initWaypointInfoListModel(m)
				return m, nil
			}
		case "u", "U":
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				// perform an unseal of waypoint, but if the waypoint is still having invalid path,
				// it will still be resealed when performing the list reinitialization
				utils.UpdateWaypointUnSeal(m.BboltDb, i.WaypointName)
				initWaypointInfoListModel(m)
				return m, nil
			}
		case "y":
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				// copy the waypoint name to the clipboard
				clipboard.WriteAll(i.WaypointName)
				return m, nil
			}
		case "Y":
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				// copy the waypoint absolute path to the clipboard
				clipboard.WriteAll(i.WaypointPath)
				return m, nil
			}
		case "j", "down":
			m.WaypointInfoList.CursorDown()
			m.WaypointInfoListCursorPosition = m.WaypointInfoList.Index()
			// refresh the help key map to reflect the sealed state of the new selection
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				m.WaypointInfoList.AdditionalShortHelpKeys = initWaypointInfoListKeyMap(i.WaypointIsSealed)
			}
		case "k", "up":
			m.WaypointInfoList.CursorUp()
			m.WaypointInfoListCursorPosition = m.WaypointInfoList.Index()
			// refresh the help key map to reflect the sealed state of the new selection
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				m.WaypointInfoList.AdditionalShortHelpKeys = initWaypointInfoListKeyMap(i.WaypointIsSealed)
			}
		case "ctrl+c", "esc", "q":
			m.IsQuit = true
			return m, tea.Quit
		}
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
	return tea.NewView(m.WaypointInfoList.View())
}
