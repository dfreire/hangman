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
	db  *bolt.DB
	app *app.App
}

func NewApp(db *bolt.DB, gormDB gorm.DB) *HangmanApp {

	app := app.NewApp(db)

	exists, err := app.ExistsApp(appId)
	if err != nil {
		panic(err)
	}
	if !exists {
		app.CreateApp(appId)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(gamesBucketName))
		return err
	})
	if err != nil {
		panic(err)
	}

	return &HangmanApp{db: db, app: app}
}
