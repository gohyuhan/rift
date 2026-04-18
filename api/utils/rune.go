package utils

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
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
func NormalizeAndCheckRuneCommandsAreValid(commandsString string) ([]*pb.RuneCmds, error) {
	var normalizedRuneCmds []*pb.RuneCmds
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
				normalizedRuneCmds = append(normalizedRuneCmds, &pb.RuneCmds{Commands: cmdArray})
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
func ParseRuneCommandsToString(runeCommands []*pb.RuneCmds) string {
	var parsedString strings.Builder
	for _, rune := range runeCommands {
		parsedString.WriteString(strings.Join(rune.Commands, " "))
	}
	return parsedString.String()
}

// ----------------------------------
//
//	Looks up the rune for waypointntPath for trigger during waypoint navigation.
//	Returns true and the deserialized *pb.Rune if a rune record exists; returns
//	false and an empty *pb.Rune on any error or if no rune is found. Never
//	returns an error — all failures are silently ignored so that waypoint path
//	changes are never blocked by rune retrieval issues.
//
// ----------------------------------
func RetrieveRuneForTrigger(waypointntPath string) (bool, *pb.Rune) {
	rune := &pb.Rune{}
	hasRuneToTrigger := false

	bboltReadDb, bboltReadDbErr := db.OpenReadDB()
	if bboltReadDbErr != nil {
		return hasRuneToTrigger, rune
	}
	defer db.CloseDB(bboltReadDb)
	_ = bboltReadDb.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(db.RuneBucket)
		if bucket == nil {
			return nil
		}

		existing := bucket.Get([]byte(waypointntPath))
		if existing == nil {
			return nil
		}

		if err := proto.Unmarshal(existing, rune); err != nil {
			return nil
		}
		hasRuneToTrigger = true
		return nil
	})

	return hasRuneToTrigger, rune
}

// ----------------------------------
//
//	Looks up the rune associated with the given waypoint path within an open
//	bbolt transaction. Returns the rune bucket and the deserialized *pb.Rune,
//	or an error if the bucket is missing, no rune exists for the path, or the
//	stored data cannot be unmarshalled.
//
// ----------------------------------
func RetrieveRuneBasedOnWaypointPath(tx *bbolt.Tx, waypointntName string) (*bbolt.Bucket, *pb.Rune, bool, error) {
	bucket := tx.Bucket(db.RuneBucket)
	if bucket == nil {
		return nil, nil, false, fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RuneBucketNotFoundError, style.ColorError, false))
	}

	existing := bucket.Get([]byte(waypointntName))
	if existing == nil {
		return bucket, &pb.Rune{EnterRunes: nil, LeaveRunes: nil}, false, nil
	}

	rune := &pb.Rune{}
	if err := proto.Unmarshal(existing, rune); err != nil {
		return nil, nil, true, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RuneDataCorruptedError, waypointntName), style.ColorError, false))
	}

	return bucket, rune, false, nil
}

// ----------------------------------
//
//	Opens a write transaction and records each path in corruptedWaypointsPath
//	into the RuneDataCorruptedBucketRecord bucket as a best-effort marker (bucket
//	errors are silently ignored). Always returns a formatted error listing all
//	corrupted paths, so the caller receives the corruption message regardless of
//	whether the write succeeded.
//
// ----------------------------------
func RecordCorruptedRuneInfo(corruptedWaypointsPath []string) error {
	// best-effort write — ignore the Update error; the caller always gets the corruption message
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)

	bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		runeCorruptedBucket := tx.Bucket(db.RuneDataCorruptedBucketRecord)
		if runeCorruptedBucket != nil {
			for _, corruptedWaypointPath := range corruptedWaypointsPath {
				runeCorruptedBucket.Put([]byte(corruptedWaypointPath), []byte(corruptedWaypointPath))
			}
		}
		return nil
	})
	return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RuneDataCorruptedError, strings.Join(corruptedWaypointsPath, ",")), style.ColorError, false))
}

// ----------------------------------
//
//	Looks up the rune for waypointPath within an existing write transaction,
//	returning the bucket and deserialized *pb.Rune ready for mutation. If no
//	rune exists for the path, returns the bucket with a new empty *pb.Rune so
//	the caller can populate and save it. Unlike RetrieveRuneBasedOnWaypointPath,
//	the corruption bool is omitted — callers that need it should use that function instead.
//
// ----------------------------------
func GetRuneForUpdate(tx *bbolt.Tx, waypointPath string) (*bbolt.Bucket, *pb.Rune, error) {
	bucket := tx.Bucket(db.RuneBucket)
	if bucket == nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RuneBucketNotFoundError, style.ColorError, false))
	}

	existing := bucket.Get([]byte(waypointPath))
	if existing == nil {
		return bucket, &pb.Rune{EnterRunes: nil, LeaveRunes: nil}, nil
	}

	rune := &pb.Rune{}
	if err := proto.Unmarshal(existing, rune); err != nil {
		return nil, nil, fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RuneDataCorruptedError, waypointPath), style.ColorError, false))
	}

	return bucket, rune, nil
}

// ----------------------------------
//
//	Marshals rune and writes it into bucket under waypointPath. Returns an error
//	if marshalling fails or the bucket write fails.
//
// ----------------------------------
func PutRune(bucket *bbolt.Bucket, waypointPath string, rune *pb.Rune) error {
	data, err := proto.Marshal(rune)
	if err != nil {
		return fmt.Errorf("%s", style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRuneUpdateError, waypointPath, err.Error()), style.ColorError, false))
	}
	return bucket.Put([]byte(waypointPath), data)
}

// ----------------------------------
//
//	Deletes the rune entry for waypointPath from both RuneBucket and
//	RuneDataCorruptedBucketRecord within the given write transaction.
//	Silently no-ops if either bucket is absent; must be called inside
//	an existing Update.
//
// ----------------------------------
func DestroyEngraveRune(tx *bbolt.Tx, waypointPath string) {
	runeBucket := tx.Bucket(db.RuneBucket)
	if runeBucket == nil {
		return
	}
	runeBucket.Delete([]byte(waypointPath))

	runeCorruptedBucket := tx.Bucket(db.RuneDataCorruptedBucketRecord)
	if runeCorruptedBucket == nil {
		return
	}
	runeCorruptedBucket.Delete([]byte(waypointPath))
}

// ----------------------------------
//
//	Clears the enter-rune slot for waypointPath in a single Update transaction.
//	If LeaveRunes is also nil after clearing, the entire rune record is deleted
//	via DestroyEngraveRune instead of writing back an empty record.
//
// ----------------------------------
func RemoveEnterRuneCmds(waypointPath string) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)
	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket, rune, err := GetRuneForUpdate(tx, waypointPath)
		if err != nil {
			return err
		}
		rune.EnterRunes = nil
		if rune.LeaveRunes == nil {
			DestroyEngraveRune(tx, waypointPath)
			return nil
		}
		return PutRune(bucket, waypointPath, rune)
	})
}

// ----------------------------------
//
//	Clears the leave-rune slot for waypointPath in a single Update transaction.
//	If EnterRunes is also nil after clearing, the entire rune record is deleted
//	via DestroyEngraveRune instead of writing back an empty record.
//
// ----------------------------------
func RemoveLeaveRuneCmds(waypointPath string) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)
	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket, rune, err := GetRuneForUpdate(tx, waypointPath)
		if err != nil {
			return err
		}
		rune.LeaveRunes = nil
		if rune.EnterRunes == nil {
			DestroyEngraveRune(tx, waypointPath)
			return nil
		}
		return PutRune(bucket, waypointPath, rune)
	})
}

// ----------------------------------
//
//	Replaces the enter-rune slot for waypointPath with EnterRunes and persists
//	the record in a single Update transaction. If no rune exists for the path,
//	a new record is created. Any previously engraved enter runes are overwritten.
//
// ----------------------------------
func EngraveEnterRuneCmds(waypointPath string, EnterRunes []*pb.RuneCmds) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)
	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket, rune, err := GetRuneForUpdate(tx, waypointPath)
		if err != nil {
			return err
		}
		rune.EnterRunes = EnterRunes
		return PutRune(bucket, waypointPath, rune)
	})
}

// ----------------------------------
//
//	Replaces the leave-rune slot for waypointPath with LeavesRunes and persists
//	the record in a single Update transaction. If no rune exists for the path,
//	a new record is created. Any previously engraved leave runes are overwritten.
//
// ----------------------------------
func EngraveLeaveRuneCmds(waypointPath string, LeavesRunes []*pb.RuneCmds) error {
	bboltWriteDb, bboltWriteDbErr := db.OpenWriteDB()
	if bboltWriteDbErr != nil {
		return bboltWriteDbErr
	}
	defer db.CloseDB(bboltWriteDb)
	return bboltWriteDb.Update(func(tx *bbolt.Tx) error {
		bucket, rune, err := GetRuneForUpdate(tx, waypointPath)
		if err != nil {
			return err
		}
		rune.LeaveRunes = LeavesRunes
		return PutRune(bucket, waypointPath, rune)
	})
}
