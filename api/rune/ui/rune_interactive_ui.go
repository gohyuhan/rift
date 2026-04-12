package ui

import (
	"fmt"
	"os"
	"sync/atomic"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
)

// ----------------------------------
//
//	bubbletea model that drives the rune engraving interactive UI;
//	holds all display and engraving state across Init / Update / View
//
// ----------------------------------
type RuneInteractiveModel struct {
	ChosenWaypointName string
	RuneEngraved       atomic.Bool
	IsQuit             bool
	ShowPopUp          atomic.Bool
	PopUpType          string
	IsTypingMode       atomic.Bool
	ErrMessage         error
	RunePopUpModel     interface{}
	Width              int
	Height             int
	IsRenderInit       atomic.Bool
	ExistingEnterRune  []*pb.Rune
	ExistingLeaveRune  []*pb.Rune
}

// ----------------------------------
//
//	Allocates and initialises a RuneInteractiveModel for the given waypoint,
//	opening the ChooseRuneEngraveOptionPopUp immediately so the user can pick
//	whether to engrave an on-enter or on-leave rune.
//
// ----------------------------------
func initRuneInteractiveModel(waypointName string) *RuneInteractiveModel {
	runeInteractiveModel := RuneInteractiveModel{
		IsQuit:             false,
		PopUpType:          ChooseRuneEngraveOptionPopUp,
		ChosenWaypointName: waypointName,
	}
	runeInteractiveModel.ShowPopUp.Store(true)
	runeInteractiveModel.IsTypingMode.Store(false)
	runeInteractiveModel.RuneEngraved.Store(false)

	return &runeInteractiveModel
}

// ----------------------------------
//
//	Runs the bubbletea program for rune engraving on the named waypoint;
//	returns true if a rune was successfully engraved, false if the user
//	cancelled. Returns an error if the TUI itself fails to start or run.
//
// ----------------------------------
func RunRuneInteractive(waypointName string) (bool, error) {
	runeInteractiveModel := initRuneInteractiveModel(waypointName)
	// route program output to stderr so stdout stays clean for callers
	p := tea.NewProgram(runeInteractiveModel, tea.WithOutput(os.Stderr))
	result, err := p.Run()
	if err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RuneInteractiveError, err), style.ColorError, false)
		return false, fmt.Errorf("%s", errorMessage)
	}

	// safe: Run always returns the model passed to NewProgram
	final := result.(*RuneInteractiveModel)
	if final.ErrMessage != nil {
		return false, final.ErrMessage
	}

	return final.RuneEngraved.Load(), nil
}

// ----------------------------------
//
//	satisfies the tea.Model interface; no startup commands are needed
//
// ----------------------------------
func (m *RuneInteractiveModel) Init() tea.Cmd {
	m.PopUpType = ChooseRuneEngraveOptionPopUp
	return initChooseRuneEngraveOptionPopUpModel(m)
}

// ----------------------------------
//
//	Satisfies the tea.Model interface; handles window resize and key input.
//	ctrl+c / q quits unconditionally; esc steps back from the rune entry
//	popup to the type-selection popup, or quits from the type-selection popup.
//
// ----------------------------------
func (m *RuneInteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.IsRenderInit.CompareAndSwap(false, true)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.IsQuit = true
			return m, tea.Quit
		case "esc":
			switch m.PopUpType {
			case ChooseRuneEngraveOptionPopUp:
				m.RuneEngraved.Store(false)
				m.IsQuit = true
				return m, tea.Quit
			case EngraveRuneCommandsPopUp:
				m.IsTypingMode.Store(false)
				m.PopUpType = ChooseRuneEngraveOptionPopUp
				cmd = initChooseRuneEngraveOptionPopUpModel(m)
				return m, cmd
			}

			return m, nil
		}

		if m.IsTypingMode.Load() {
			m, cmd = handleTypingInteraction(m, msg)
		} else {
			m, cmd = handleNonTypingInteraction(m, msg)
		}

		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	}
	return m, tea.Batch(cmds...)
}

// ----------------------------------
//
//	Satisfies the tea.Model interface; renders the rune engraving UI.
//	Returns an empty view once the user has quit to avoid a final flash
//	of stale content.
//
// ----------------------------------
func (m *RuneInteractiveModel) View() tea.View {
	if m.IsQuit {
		return tea.NewView("")
	}
	return tea.NewView(renderRuneInteractiveUIView(m))
}
