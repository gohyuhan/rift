package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/rift/api/spell"
)

// ----------------------------------
//
//	processes key events when the model is not in typing mode;
//	j/↓ and k/↑ move the cursor; when the help popup is open, j/↓ and k/↑
//	scroll the viewport instead; ? opens the help popup; enter casts the
//	selected spell; backspace forgets the selected spell and rebuilds the list
//
// ----------------------------------
func handleNonTypingInteraction(m *SpellbookInteractiveModel, msg tea.KeyPressMsg) (*SpellbookInteractiveModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			m.SpellHelpViewport, cmd = m.SpellHelpViewport.Update(msg)
			return m, cmd
		}
		m.SpellInfoList.CursorDown()
		m.SpellInfoListCursorPosition = m.SpellInfoList.Index()
	case "k", "up":
		if m.ShowPopUp.Load() {
			var cmd tea.Cmd
			m.SpellHelpViewport, cmd = m.SpellHelpViewport.Update(msg)
			return m, cmd
		}
		m.SpellInfoList.CursorUp()
		m.SpellInfoListCursorPosition = m.SpellInfoList.Index()
	case "?":
		m.ShowPopUp.Store(true)
		m.PopUpType = HelpPopUp
	case "enter":
		if _, ok := m.SpellInfoList.SelectedItem().(spellInfoItem); ok {
			return m, nil
		}
	case "backspace":
		if i, ok := m.SpellInfoList.SelectedItem().(spellInfoItem); ok {
			// perform a forget on the selected spell
			forgetErr := spell.ForgetSpell(m.BboltDb, i.SpellName, false)
			if forgetErr != nil {
				m.ErrMessage = forgetErr
				m.IsQuit = true
				return m, tea.Quit
			}
			initSpellInfoListModel(m)
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
func handleTypingInteraction(m *SpellbookInteractiveModel, msg tea.KeyPressMsg) (*SpellbookInteractiveModel, tea.Cmd) {
	switch msg.String() {
	}

	return m, nil
}
