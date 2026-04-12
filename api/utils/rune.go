package utils

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"mvdan.cc/sh/v3/shell"
)

// ----------------------------------
//
//	Splits commandsString by newline, shell-parses each non-empty line, and
//	rejects any line whose first token is "cd" (path navigation is forbidden
//	in rune commands). Returns the parsed commands as a slice of *pb.Rune
//	(one per non-empty line), or an error if any line fails shell parsing or
//	contains a cd command. Blank lines are silently skipped.
//
// ----------------------------------
func NormalizeAndCheckRuneCommandsAreValid(commandsString string) ([]*pb.Rune, error) {
	var normalizedRuneCmds []*pb.Rune
	cmdsStringArray := strings.Split(commandsString, "\n")
	for _, cmdString := range cmdsStringArray {
		cmdString = strings.TrimSpace(cmdString)
		cmdArray, cmdArrayErr := shell.Fields(cmdString, nil)
		if cmdArrayErr != nil {
			return nil, cmdArrayErr
		}
		if len(cmdArray) > 0 {
			if cmdArray[0] == "cd" {
				return nil, fmt.Errorf("%s", i18n.LANGUAGEMAPPING.RuneCommandsInvalidDueToCDCommand)
			} else {
				normalizedRuneCmds = append(normalizedRuneCmds, &pb.Rune{Commands: cmdArray})
			}
		}
	}

	return normalizedRuneCmds, nil
}

// ----------------------------------
//
//	Joins a slice of *pb.Rune command arrays into a single string suitable for
//	pre-populating the rune commands textarea. Each Rune's tokens are space-joined
//	and the results are concatenated in order. Returns an empty string when
//	runeCommands is nil or empty.
//
// ----------------------------------
func ParseRuneCommandsToString(runeCommands []*pb.Rune) string {
	var parsedString strings.Builder
	for _, rune := range runeCommands {
		parsedString.WriteString(strings.Join(rune.Commands, " "))
	}
	return parsedString.String()
}
