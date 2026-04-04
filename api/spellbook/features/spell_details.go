package features

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
//	Fetches a single named spell and builds a detail display list.
//	Uses a read-only View transaction; corruption recording is deferred to a
//	separate Update transaction after View completes (same pattern as
//	RetrieveWaypointInfoDetail).
//	The returned slice has one labelled row per field.
//	Labels are padded to uniform width via style.PadAndRenderLabels so values
//	align in a clean column regardless of the active language.
//
// ----------------------------------
func RetrieveSpellInfoDetail(bboltDb *bbolt.DB, spellName string) ([]string, error) {
	var waypointDetailInfo []string
	waypointCorrupted := false

	viewErr := bboltDb.View(func(tx *bbolt.Tx) error {
		// ensure the spell bucket exists before looking up the key
		spellBucket := tx.Bucket(db.SpellBucket)
		if spellBucket == nil {
			return fmt.Errorf("%s", style.RenderStringWithColor(i18n.LANGUAGEMAPPING.SpellBucketNotFoundError, style.ColorError, false))
		}

		// look up the specific spell by name
		existing := spellBucket.Get([]byte(spellName))
		if existing == nil {
			errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.RiftSpellDoNotExistsError, spellName), style.ColorError, false)
			return fmt.Errorf("%s", errorMessage)
		}

		// deserialize the stored proto; flag corruption and exit the View cleanly —
		// the actual corruption recording write is deferred to after the View completes
		existingSpell := &pb.Spell{}
		protoErr := proto.Unmarshal(existing, existingSpell)
		if protoErr != nil {
			waypointCorrupted = true
			return nil
		}

		// collect raw labels and pad them to uniform width
		rawLabels := []string{
			i18n.LANGUAGEMAPPING.RiftSpellDetailName,
			i18n.LANGUAGEMAPPING.RiftSpellDetailCommand,
			i18n.LANGUAGEMAPPING.RiftSpellDetailAddedAt,
			i18n.LANGUAGEMAPPING.RiftSpellDetailCastCount,
		}
		paddedLabels := style.PadAndRenderLabels(rawLabels, style.ColorBlueGrayMuted, true)

		// --- spell name (highlight in cyan; muted colors are reserved for sealed waypoints) ---
		nameValue := style.RenderStringWithColor(spellName, style.ColorCyanSoft, false)
		waypointDetailInfo = append(waypointDetailInfo, paddedLabels[0]+"  "+nameValue)

		// --- spell cmd ---
		cmdValue := style.RenderStringWithColor(strings.Join(existingSpell.SpellCommand, " "), style.ColorBlueMuted, false)
		waypointDetailInfo = append(waypointDetailInfo, paddedLabels[1]+"  "+cmdValue)

		// --- spell added at (stored as RFC3339 UTC, displayed in local time) ---
		addAtDisplay := existingSpell.SpellAddedAt
		if parsed, err := time.Parse(time.RFC3339, existingSpell.SpellAddedAt); err == nil {
			addAtDisplay = parsed.Local().Format("2006-01-02 15:04:05")
		}
		addAtValue := style.RenderStringWithColor(addAtDisplay, style.ColorBlueMuted, false)
		waypointDetailInfo = append(waypointDetailInfo, paddedLabels[2]+"  "+addAtValue)

		// --- spell cast count ---
		castValue := style.RenderStringWithColor(strconv.FormatInt(existingSpell.SpellCastCount, 10), style.ColorBlueMuted, false)
		waypointDetailInfo = append(waypointDetailInfo, paddedLabels[3]+"  "+castValue)

		return nil
	})

	// View is complete — safe to open a separate Update for the corruption write
	if waypointCorrupted {
		viewErr = apiUtils.RecordCorruptedSpellInfo(bboltDb, []string{spellName})
	}

	return waypointDetailInfo, viewErr
}
