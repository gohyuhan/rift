package ui

// ListItemOrTitleWidthPad is the number of columns reserved for list chrome
// (borders, cursor prefix, padding) so that text can be safely truncated to
// the remaining available width without overflowing the terminal.
const (
	ListItemOrTitleWidthPad    = 4
	TextInputWidthPad          = 6
	ListTitleHorizontalPadding = 2
)

// PopUpType sentinel values stored on SpellbookInteractiveModel.PopUpType to
// indicate which popup (if any) is currently active.
const (
	ChooseRuneEngraveTypePopUp = "ChooseRuneEngraveTypePopUp"
	EnterRunePopUp             = "EnterRunePopUp"
)
