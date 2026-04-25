package utils

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gohyuhan/rift/constant"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
	"mvdan.cc/sh/v3/shell"
	"mvdan.cc/sh/v3/syntax"
)

// ----------------------------------
//
//	Splits commandsString by newline, shell-parses each non-empty line, and
//	rejects any line whose first token is a shell built-in (e.g. cd, export,
//	source) or a rift navigation command (path navigation is forbidden in
//	ritual commands). Returns the parsed commands as a slice of *pb.RitualCmds
//	(one per non-empty line), or an error if any line fails shell parsing or
//	contains a forbidden command. Blank lines are silently skipped.
//
// ----------------------------------
func NormalizeAndCheckRitualCommandsAreValid(commandsString string) ([]*pb.RitualCmds, error) {
	var normalizedRitualCmds []*pb.RitualCmds
	cmdsStringArray := strings.Split(commandsString, "\n")
	for _, cmdString := range cmdsStringArray {
		cmdString = strings.TrimSpace(cmdString)
		cmdArray, cmdArrayErr := shell.Fields(cmdString, nil)
		if cmdArrayErr != nil {
			return nil, cmdArrayErr
		}
		if len(cmdArray) > 0 {
			if slices.Contains(constant.ShellBuildInCmd, cmdArray[0]) {
				msg := fmt.Sprintf(i18n.LANGUAGEMAPPING.RitualCommandsInvalidDueToShellBuildInCommand, utils.ShellBuiltinExample())
				return nil, fmt.Errorf("%s", msg)
			} else if utils.IsRiftNavigationCommand(cmdArray) {
				errMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.ForbiddenRiftNavigationRitualCommand, style.ColorError, false)
				return nil, fmt.Errorf("%s", errMessage)
			} else {
				normalizedRitualCmds = append(normalizedRitualCmds, &pb.RitualCmds{Commands: cmdArray})
			}
		}
	}

	return normalizedRitualCmds, nil
}

// ----------------------------------
//
//	Joins a slice of *pb.RitualCmds command arrays into a single string suitable
//	for pre-populating the ritual commands textarea. Each token is shell-quoted
//	if it contains characters that shell.Fields would interpret as operators or
//	word-splitters, so that the reconstructed line round-trips through
//	shell.Fields without error. Returns an empty string when ritualCommands is
//	nil or empty.
//
// ----------------------------------
func ParseRitualCommandsToString(ritualCommands []*pb.RitualCmds) string {
	var parsedString strings.Builder
	for _, ritual := range ritualCommands {
		quoted := make([]string, len(ritual.Commands))
		for i, token := range ritual.Commands {
			quoted[i], _ = syntax.Quote(token, syntax.LangBash)
		}
		parsedString.WriteString(strings.Join(quoted, " "))
		parsedString.WriteRune('\n')
	}
	return parsedString.String()
}

// ----------------------------------
//
//	Persists the ritual names into the corrupted-records bucket via a fresh
//	Update transaction (independent of any prior transaction that may have
//	failed), then returns a user-facing corruption error.
//	The Update's own error is intentionally not propagated — recording is
//	best-effort; the corruption error is always returned to the caller.
//
// ----------------------------------
func RecordCorruptedRitualInfo(corruptedRitualsName []string) error {
	// best-effort write — ignore the Update error; the caller always gets the corruption message
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		ritualCorruptedBucket := tx.Bucket(db.RitualDataCorruptedBucketRecord)
		if ritualCorruptedBucket != nil {
			for _, corruptedRitual := range corruptedRitualsName {
				ritualCorruptedBucket.Put([]byte(corruptedRitual), []byte(corruptedRitual))
			}
		}
		return nil
	})
	return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RitualDataCorruptedError, strings.Join(corruptedRitualsName, ",")), style.ColorError, false))
}

// ----------------------------------
//
//	Fetches and deserializes the named ritual from the bucket within an
//	already-open Update transaction. Returns the bucket, the deserialized
//	record, or an error if the bucket is missing, the ritual does not exist,
//	or the stored proto is corrupted. Callers mutate the returned record and
//	re-persist it via bucket.Put.
//
// ----------------------------------
func GetRitualForUpdate(tx *bbolt.Tx, ritualName string) (*bbolt.Bucket, *pb.Ritual, error) {
	bucket := tx.Bucket(db.RitualBucket)
	if bucket == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RitualBucketNotFoundError, style.ColorError, false))
	}

	existing := bucket.Get([]byte(ritualName))
	if existing == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualDoNotExistsError, ritualName), style.ColorError, false))
	}

	ritual := &pb.Ritual{}
	if err := proto.Unmarshal(existing, ritual); err != nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RitualDataCorruptedError, ritualName), style.ColorError, false))
	}

	return bucket, ritual, nil
}

// ----------------------------------
//
//	Persists a mutated ritual record back into its bucket. Returns an error
//	if marshalling or the bucket write fails.
//
// ----------------------------------
func PutRitual(bucket *bbolt.Bucket, ritualName string, ritual *pb.Ritual) error {
	data, err := proto.Marshal(ritual)
	if err != nil {
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualUpdateError, ritualName, err.Error()), style.ColorError, false))
	}
	return bucket.Put([]byte(ritualName), data)
}

// ----------------------------------
//
//	Increments the invoked count for the named ritual in the DB bucket.
//	Returns an error if the bucket is missing, the data is unreadable, or the
//	ritual does not exist.
//
// ----------------------------------
func UpdateRitualInvokedCount(ritualName string) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket, ritual, err := GetRitualForUpdate(tx, ritualName)
		if err != nil {
			return err
		}
		ritual.RitualInvokeCount += 1
		return PutRitual(bucket, ritualName, ritual)
	})
}
