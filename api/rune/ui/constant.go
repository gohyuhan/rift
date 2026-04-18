package ui

// ListItemOrTitleWidthPad is the number of columns reserved for list chrome
// (borders, cursor prefix, padding) so that text can be safely truncated to
// the remaining available width without overflowing the terminal.
const (
	ListItemOrTitleWidthPad    = 4
	TextInputWidthPad          = 6
	ListTitleHorizontalPadding = 2
)

// PopUpType sentinel values stored on RuneInteractiveModel.PopUpType to
// indicate which popup (if any) is currently active.
const (
	ChooseRuneEngraveOptionPopUp = "ChooseRuneEngraveOptionPopUp"
	EngraveRuneCommandsPopUp     = "EngraveRuneCommandsPopUp"
)

// Engrave/remove option types selected from the ChooseRuneEngraveOptionPopUp;
// carried into initEngraveRuneCommandsPopUpModel and the engraving write calls
// to identify whether the enter or leave rune slot is being modified.
const (
	EngraveRuneEnterType = "EngraveRuneEnterType"
	EngraveRuneLeaveType = "EngraveRuneLeaveType"
	RemoveRuneEnterType  = "RemoveRuneEnterType"
	RemoveRuneLeaveType  = "RemoveRuneLeaveType"
)
