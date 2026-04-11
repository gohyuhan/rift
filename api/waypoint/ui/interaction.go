package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/api/waypoint/features"
)

// ----------------------------------
//
//	processes key events when the model is not in typing mode;
//	j/↓ and k/↑ move the cursor; when the help popup is open, j/↓ and k/↑
//	scroll the viewport instead; u/U attempts to unseal and rebuilds the list;
//	r opens the rebind path input popup; R opens the reforge name input popup;
//	y copies the waypoint name to the clipboard; Y copies the path;
//	? opens the help popup; enter navigates to the selected waypoint (skipped
//	if sealed); backspace destroys the selected waypoint and rebuilds the list
//
// ----------------------------------
func handleNonTypingInteraction(m *WaypointInteractiveModel, msg tea.KeyPressMsg) (*WaypointInteractiveModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			m.WaypointHelpViewport, cmd = m.WaypointHelpViewport.Update(msg)
			return m, cmd
		}
		m.WaypointInfoList.CursorDown()
		m.WaypointInfoListCursorPosition = m.WaypointInfoList.Index()
	case "k", "up":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			m.WaypointHelpViewport, cmd = m.WaypointHelpViewport.Update(msg)
			return m, cmd
		}
		m.WaypointInfoList.CursorUp()
		m.WaypointInfoListCursorPosition = m.WaypointInfoList.Index()
	case "r":
		// open the rebind path input popup for the selected waypoint
		selectedWaypoint := m.WaypointInfoList.SelectedItem()
		if selectedWaypoint != nil {
			parsedWaypointName := selectedWaypoint.(waypointInfoItem).WaypointName
			m.ShowPopUp.Store(true)
			m.IsTypingMode.Store(true)
			m.PopUpType = RebindPopUp
			initRebindPopUp(m, parsedWaypointName)
		}
		return m, nil
	case "R":
		// open the reforge name input popup for the selected waypoint
		selectedWaypoint := m.WaypointInfoList.SelectedItem()
		if selectedWaypoint != nil {
			parsedWaypointName := selectedWaypoint.(waypointInfoItem).WaypointName
			m.ShowPopUp.Store(true)
			m.IsTypingMode.Store(true)
			m.PopUpType = ReforgePopUp
			initReforgePopUp(m, parsedWaypointName)
		}
		return m, nil
	case "u", "U":
		if i, ok := m.WaypointInfoList.SelectedItem().(waypointInfoItem); ok {
			// perform an unseal of waypoint, but if the waypoint is still having invalid path,
			// it will still be resealed when performing the list reinitialization
			utils.UpdateWaypointUnSeal(i.WaypointName)
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
	case "?":
		m.ShowPopUp.Store(true)
		m.PopUpType = HelpPopUp
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
			destoryErr := features.DestroyDiscoveredWaypoint(i.WaypointName, false)
			if destoryErr != nil {
				m.ErrMessage = destoryErr
				m.IsQuit = true
				return m, tea.Quit
			}
			initWaypointInfoListModel(m)
			return m, nil
		}
	}
	return m, nil
}

// ----------------------------------
//
//	processes key events when the model is in typing mode (e.g. an input popup
//	is active); on enter, submits the active popup's input: for RebindPopUp and
//	ReforgePopUp it calls the stored trigger function and closes the popup on
//	success, or surfaces the error inline on failure; for all other keys,
//	forwards the event to the active popup's input component so text editing
//	works normally
//
// ----------------------------------
func handleTypingInteraction(m *WaypointInteractiveModel, msg tea.KeyPressMsg) (*WaypointInteractiveModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+y":
		switch m.PopUpType {
		case RebindPopUp:
			popUp, ok := m.WaypointPopUpModel.(*RebindPopUpModel)
			if ok {
				clipboard.WriteAll(popUp.RebindPathInput.Value())
			}
		case ReforgePopUp:
			popUp, ok := m.WaypointPopUpModel.(*ReforgePopUpModel)
			if ok {
				clipboard.WriteAll(popUp.ReforgeWaypointNameInput.Value())
			}
		}

	case "ctrl+p":
		content, err := clipboard.ReadAll()
		if err != nil {
			return m, nil
		}
		msg := tea.PasteMsg{
			Content: content,
		}
		switch m.PopUpType {
		case RebindPopUp:
			popUp, ok := m.WaypointPopUpModel.(*RebindPopUpModel)
			if ok {
				popUp.RebindPathInput, cmd = popUp.RebindPathInput.Update(msg)
			}
		case ReforgePopUp:
			popUp, ok := m.WaypointPopUpModel.(*ReforgePopUpModel)
			if ok {
				popUp.ReforgeWaypointNameInput, cmd = popUp.ReforgeWaypointNameInput.Update(msg)
			}
		}
		return m, cmd

	case "enter":
		closePopUp := false

		switch m.PopUpType {
		case RebindPopUp:
			popUp, ok := m.WaypointPopUpModel.(*RebindPopUpModel)
			if ok {
				newWaypointPath := strings.TrimSpace(popUp.RebindPathInput.Value())
				rebindErr := popUp.OnInputFuncTrigger(popUp.WaypointName, newWaypointPath, false)
				if rebindErr != nil {
					// surface the error in the popup rather than quitting
					popUp.Error = rebindErr
					return m, nil
				}
				closePopUp = true
				initWaypointInfoListModel(m)
			}
		case ReforgePopUp:
			popUp, ok := m.WaypointPopUpModel.(*ReforgePopUpModel)
			if ok {
				newWaypointName := strings.TrimSpace(popUp.ReforgeWaypointNameInput.Value())
				reforgeErr := popUp.OnInputFuncTrigger(popUp.WaypointName, newWaypointName, false)
				if reforgeErr != nil {
					// surface the error in the popup rather than quitting
					popUp.Error = reforgeErr
					return m, nil
				}
				closePopUp = true
				initWaypointInfoListModel(m)
			}
		}

		if closePopUp {
			m.ShowPopUp.Store(false)
			m.IsTypingMode.Store(false)
			m.PopUpType = NoPopUp
		}

		return m, nil
	}

	// forward all other key events to the active popup's input component
	switch m.PopUpType {
	case RebindPopUp:
		popUp, ok := m.WaypointPopUpModel.(*RebindPopUpModel)
		if ok {
			popUp.RebindPathInput, cmd = popUp.RebindPathInput.Update(msg)
			return m, cmd
		}
	case ReforgePopUp:
		popUp, ok := m.WaypointPopUpModel.(*ReforgePopUpModel)
		if ok {
			popUp.ReforgeWaypointNameInput, cmd = popUp.ReforgeWaypointNameInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}
