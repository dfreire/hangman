package hangman

import (
	"github.com/boltdb/bolt"
	"github.com/dfreire/hangman/deltas"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/puffinframework/app"
)

const (
	appId = "HangmanApp"
)

type HangmanApp struct {
	gormDB       gorm.DB
	app          *app.App
	deltaService deltas.DeltaService
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

	gormDB.CreateTable(Game{})

	return &HangmanApp{
		gormDB:       gormDB,
		app:          app,
		deltaService: deltas.NewBoltDeltaService(boltDB, "HangmanDeltas"),
	}
}
