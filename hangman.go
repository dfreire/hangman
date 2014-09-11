package hangman

import (
	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/puffinframework/app"
)

const (
	appId = "HangmanApp"
)

type HangmanApp struct {
	boltDB *bolt.DB
	gormDB gorm.DB
	app    *app.App
}

func NewApp(boltDB *bolt.DB, gormDB gorm.DB) *HangmanApp {

	app := app.NewApp(boltDB)

	exists, err := app.ExistsApp(appId)
	if err != nil {
		panic(err)
	}
	if !exists {
		app.CreateApp(appId)
	}

	boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(gamesBucketName))
		return err
	})
	if err != nil {
		panic(err)
	}

	gormDB.CreateTable(Game{})

	return &HangmanApp{boltDB: boltDB, gormDB: gormDB, app: app}
}
