package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gohyuhan/rift/i18n"
	"github.com/gohyuhan/rift/style"
	"github.com/gohyuhan/rift/utils"
	"go.etcd.io/bbolt"
)

var (
	WaypointBucket                    = []byte("rift-waypoint")
	WaypointDataCorruptedBucketRecord = []byte("rift-waypoint-data-corrupted")
	SpellBucket                       = []byte("rift-spell")
	SpellDataCorruptedBucketRecord    = []byte("rift-spell-data-corrupted")
	RuneBucket                        = []byte("rift-rune")
	RuneDataCorruptedBucketRecord     = []byte("rift-rune-data-corrupted")
	RitualBucket                      = []byte("rift-ritual")
	RitualDataCorruptedBucketRecord   = []byte("rift-ritual-data-corrupted")
)

// ----------------------------------
//
//	Ensures the DB directory exists, opens the database, and initializes
//	all required buckets. Called once during setup.
//
// ----------------------------------
func SetupDB() error {
	dbPath, dbPathErr := utils.GetRiftDBFilePath()
	if dbPathErr != nil {
		return dbPathErr
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.DBSetupError, err.Error()), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}

	bucketSetupErr := SetupBuckets()
	if bucketSetupErr != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.DBSetupError, bucketSetupErr.Error()), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

// ----------------------------------
//
//	Creates the waypoint and spell buckets (and their corrupted-record
//	counterparts) if they do not already exist.
//
// ----------------------------------
func SetupBuckets() error {
	bboltWriteDB, bboltWriteDBErr := OpenWriteDB()
	if bboltWriteDBErr != nil {
		return bboltWriteDBErr
	}
	defer CloseDB(bboltWriteDB)
	return bboltWriteDB.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(WaypointBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(WaypointDataCorruptedBucketRecord); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(SpellBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(SpellDataCorruptedBucketRecord); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(RuneBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(RuneDataCorruptedBucketRecord); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(RitualBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(RitualDataCorruptedBucketRecord); err != nil {
			return err
		}
		return nil
	})
}

// ----------------------------------
//
//	Opens the bbolt database at the resolved DB file path with a 8-second
//	timeout. Opens read-only when isWrite is false, read-write otherwise,
//	to avoid indefinite blocking if the file is already locked.
//
// ----------------------------------
func openDB(isWrite bool) (*bbolt.DB, error) {
	dbPath, dbPathErr := utils.GetRiftDBFilePath()
	if dbPathErr != nil {
		return nil, dbPathErr
	}

	maxDBWaitSeconds := 8 * time.Second

	dbOption := &bbolt.Options{ReadOnly: true, Timeout: maxDBWaitSeconds}
	if isWrite {
		dbOption = &bbolt.Options{Timeout: maxDBWaitSeconds}
	}

	db, err := bbolt.Open(dbPath, 0o600, dbOption)
	if err != nil {
		errorMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.DBOpenError, style.ColorError, false)
		return nil, fmt.Errorf("%s", errorMessage)
	}
	return db, nil
}

// ----------------------------------
//
//	Opens the bbolt database in read-only mode.
//
// ----------------------------------
func OpenReadDB() (*bbolt.DB, error) {
	return openDB(false)
}

// ----------------------------------
//
//	Opens the bbolt database in read-write mode.
//
// ----------------------------------
func OpenWriteDB() (*bbolt.DB, error) {
	return openDB(true)
}

// ----------------------------------
//
//	Closes the bbolt database. Intended to be called via defer after OpenReadDB or OpenWriteDB.
//
// ----------------------------------
func CloseDB(db *bbolt.DB) {
	db.Close()
}
