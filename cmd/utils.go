package cmd

import (
	"fmt"
	"slices"

	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
)

var ReservedCommandKeywords = []string{
	"rift",
	"awaken",
	"discover",
	"spell",
	"cast",
	"ritual",
	"summon",
	"deploy",
	"rune",
	"seer",
	"recall",
	"loot",
	"waypoint",
	"grimore",
	"lore",
	"stats",
}

// ----------------------------------
//
// This is to check if those waypoint name defined by the user didn't conflict with rift's reserved keyword,
// such as `awaken`.
//
// ----------------------------------
func CheckIfKeywordIsReservedForRift(arg string) error {
	if slices.Contains(ReservedCommandKeywords, arg) {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftReservedKeywordError, arg), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}
