package ui

import (
	"fmt"
	"os"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
)

type WaypointInteractiveModel struct {
	SelectedWaypointPath           string
	SelectedWaypointName           string
	IsQuit                         bool
	ErrMessage                     string
	WaypointInfoList               list.Model
	WaypointInfoListCursorPosition int
	Width                          int
	Height                         int
	BboltDb                        *bbolt.DB
	IsRenderInit                   atomic.Bool
}

type waypointInfo struct {
	WaypointName         string
	WaypointPath         string
	WaypointIsSealed     bool
	WaypointSealedReason string
}

func initWaypointInteractiveModel(bboltDb *bbolt.DB) *WaypointInteractiveModel {
	waypointInteractiveModel := WaypointInteractiveModel{
		SelectedWaypointPath:           "",
		IsQuit:                         false,
		ErrMessage:                     "",
		WaypointInfoListCursorPosition: 0,
		BboltDb:                        bboltDb,
	}

	return &waypointInteractiveModel
}

func RunWaypointInteractive(bboltDb *bbolt.DB) (string, string, error) {
	waypointInteractiveModel := initWaypointInteractiveModel(bboltDb)
	p := tea.NewProgram(waypointInteractiveModel, tea.WithOutput(os.Stderr))
	result, err := p.Run()
	if err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.WaypointInteractiveError, err), style.ColorError, false)
		return "", "", fmt.Errorf("%s", errorMessage)
	}

	final := result.(*WaypointInteractiveModel)
	if final.SelectedWaypointPath == "" || final.SelectedWaypointName == "" {
		return "", "", fmt.Errorf(i18n.LANGUAGEMAPPING.RiftWaypointPathEmptyError, final.SelectedWaypointName)
	}
	return final.SelectedWaypointPath, final.SelectedWaypointName, nil
}

func (m *WaypointInteractiveModel) Init() tea.Cmd {
	return nil
}

func (m *WaypointInteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// Initialize list components once, immediately after the first window resize.
		// Valid dimensions are required to calculate item layouts (specifically text truncation);
		// initializing earlier would cause the UI layout to break.
		if m.IsRenderInit.CompareAndSwap(false, true) {
			if err := initWaypointInfoListModel(m); err != nil {
				m.ErrMessage = err.Error()
			}
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				if !i.WaypointIsSealed {
					m.SelectedWaypointPath = i.WaypointPath
					m.SelectedWaypointName = i.WaypointName
					m.IsQuit = true
					return m, tea.Quit
				}
				return m, nil
			}
		case "j", "down":
			m.WaypointInfoList.CursorDown()
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				initWaypointInfoListKeyMap(i.WaypointIsSealed)
			}
		case "k", "up":
			m.WaypointInfoList.CursorUp()
			if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
				initWaypointInfoListKeyMap(i.WaypointIsSealed)
			}
		case "ctrl+c", "esc", "q":
			m.IsQuit = true
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *WaypointInteractiveModel) View() tea.View {
	if m.IsQuit {
		return tea.NewView("")
	}
	return tea.NewView(m.WaypointInfoList.View())
}
