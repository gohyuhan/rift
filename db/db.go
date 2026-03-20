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
	SettingsBucket = []byte("rift-settings")
	WaypointBucket = []byte("rift-waypoint")
)

func SetupDB() error {
	dbPath, dbPathErr := utils.GetRiftDBFilePath()
	if dbPathErr != nil {
		return dbPathErr
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.DBSetupError, err.Error()), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}

	bboltDB, bboltDBErr := OpenDB()
	if bboltDBErr != nil {
		return bboltDBErr
	}
	defer CloseDB(bboltDB)

	bucketSetupErr := SetupBuckets(bboltDB)
	if bucketSetupErr != nil {
		errorMessage := style.RenderStringWithColor(fmt.Sprintf(i18n.LANGUAGEMAPPING.DBSetupError, bucketSetupErr.Error()), style.ColorError, false)
		return fmt.Errorf("%s", errorMessage)
	}
	return nil
}

func SetupBuckets(db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(SettingsBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(WaypointBucket); err != nil {
			return err
		}
		return nil
	})
}

func OpenDB() (*bbolt.DB, error) {
	dbPath, dbPathErr := utils.GetRiftDBFilePath()
	if dbPathErr != nil {
		return nil, dbPathErr
	}
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		errorMessage := style.RenderStringWithColor(i18n.LANGUAGEMAPPING.DBOpenError, style.ColorError, false)
		return nil, fmt.Errorf("%s", errorMessage)
	}
	return db, nil
}

func CloseDB(db *bbolt.DB) {
	db.Close()
}
