package features

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	apiUtils "github.com/gohyuhan/rift/api/utils"
	"github.com/gohyuhan/rift/db"
	"github.com/gohyuhan/rift/i18n"
	pb "github.com/gohyuhan/rift/proto"
	"github.com/gohyuhan/rift/style"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// ----------------------------------
//
//	Fetches a single named ritual and builds a detail display list.
//	Uses a read-only View transaction; corruption recording is deferred to a
//	separate Update transaction after View completes (same pattern as
//	RetrieveWaypointInfoDetail).
//	The returned slice has one labelled row per field.
//	Labels are padded to uniform width via style.PadAndRenderLabels so values
//	align in a clean column regardless of the active language.
//
// ----------------------------------
func RetrieveRitualInfoDetail(ritualName string) ([]string, error) {
	var ritualDetailInfo []string
	ritualCorrupted := false

	// open DB so we can read ritual records
	viewErr := func() error {
		bboltReadDb, bboltReadDbErr := db.OpenReadDB()
		if bboltReadDbErr != nil {
			return bboltReadDbErr
		}
		defer db.CloseDB(bboltReadDb)

		return bboltReadDb.View(func(tx *bbolt.Tx) error {
			// ensure the ritual bucket exists before looking up the key
			ritualBucket := tx.Bucket(db.RitualBucket)
			if ritualBucket == nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RitualBucketNotFoundError, style.ColorError, false))
			}

			// look up the specific ritual by name
			existing := ritualBucket.Get([]byte(ritualName))
			if existing == nil {
				errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualDoNotExistsError, ritualName), style.ColorError, false)
				return fmt.Errorf("%s", errorMessage)
			}

			// deserialize the stored proto; flag corruption and exit the View cleanly —
			// the actual corruption recording write is deferred to after the View completes
			existingRitual := &pb.Ritual{}
			protoErr := proto.Unmarshal(existing, existingRitual)
			if protoErr != nil {
				ritualCorrupted = true
				return nil
			}

			// collect raw labels and pad them to uniform width
			rawLabels := []string{
				i18n.LANGUAGEMAPPING.RiftRitualDetailName,
				i18n.LANGUAGEMAPPING.RiftRitualDetailCommands,
				i18n.LANGUAGEMAPPING.RiftRitualDetailAddedAt,
				i18n.LANGUAGEMAPPING.RiftRitualDetailInvokeCount,
			}
			paddedLabels := style.PadAndRenderLabels(rawLabels, style.ColorBlueGrayMuted, true)

			// --- ritual name (highlight in cyan; muted colors are reserved for sealed rituals) ---
			nameValue := style.RenderStringWithColor(ritualName, style.ColorCyanSoft, false)
			ritualDetailInfo = append(ritualDetailInfo, paddedLabels[0]+"  "+nameValue)

			// --- ritual commands ---
			normalizedCommands := strings.Split(apiUtils.ParseRitualCommandsToString(existingRitual.RitualCmds), "\n")
			for i, cmd := range normalizedCommands {
				normalizedCommands[i] = style.RenderStringWithColor(cmd, style.ColorBlueMuted, false)
				if i == 0 {
					ritualDetailInfo = append(ritualDetailInfo, paddedLabels[1]+"  "+normalizedCommands[i])
				} else {
					ritualDetailInfo = append(ritualDetailInfo, strings.Repeat(" ", lipgloss.Width(paddedLabels[1]))+"  "+normalizedCommands[i])
				}
			}

			// --- ritual added at (stored as RFC3339 UTC, displayed in local time) ---
			addAtDisplay := existingRitual.RitualAddedAt
			if parsed, err := time.Parse(time.RFC3339, existingRitual.RitualAddedAt); err == nil {
				addAtDisplay = parsed.Local().Format("2006-01-02 15:04:05")
			}
			addAtValue := style.RenderStringWithColor(addAtDisplay, style.ColorBlueMuted, false)
			ritualDetailInfo = append(ritualDetailInfo, paddedLabels[2]+"  "+addAtValue)

			// --- ritual invoke count ---
			invokeValue := style.RenderStringWithColor(strconv.FormatInt(existingRitual.RitualInvokeCount, 10), style.ColorBlueMuted, false)
			ritualDetailInfo = append(ritualDetailInfo, paddedLabels[3]+"  "+invokeValue)

			return nil
		})
	}()

	// View is complete — safe to open a separate Update for the corruption write
	if ritualCorrupted {
		viewErr = apiUtils.RecordCorruptedRitualInfo([]string{ritualName})
	}

	return ritualDetailInfo, viewErr
}

func RetrieveRitualInfoDetailForEdit(ritualName string) (*pb.Ritual, error) {
	var ritualInfoForEdit *pb.Ritual
	ritualCorrupted := false

	// open DB so we can read ritual records
	viewErr := func() error {
		bboltReadDb, bboltReadDbErr := db.OpenReadDB()
		if bboltReadDbErr != nil {
			return bboltReadDbErr
		}
		defer db.CloseDB(bboltReadDb)

		return bboltReadDb.View(func(tx *bbolt.Tx) error {
			// ensure the ritual bucket exists before looking up the key
			ritualBucket := tx.Bucket(db.RitualBucket)
			if ritualBucket == nil {
				return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.RitualBucketNotFoundError, style.ColorError, false))
			}

			// look up the specific ritual by name
			existing := ritualBucket.Get([]byte(ritualName))
			if existing == nil {
				errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftRitualDoNotExistsError, ritualName), style.ColorError, false)
				return fmt.Errorf("%s", errorMessage)
			}

			// deserialize the stored proto; flag corruption and exit the View cleanly —
			// the actual corruption recording write is deferred to after the View completes
			existingRitual := &pb.Ritual{}
			protoErr := proto.Unmarshal(existing, existingRitual)
			if protoErr != nil {
				ritualCorrupted = true
				return nil
			}

			ritualInfoForEdit = existingRitual

			return nil
		})
	}()

	// View is complete — safe to open a separate Update for the corruption write
	if ritualCorrupted {
		viewErr = apiUtils.RecordCorruptedRitualInfo([]string{ritualName})
	}

	return ritualInfoForEdit, viewErr
}
