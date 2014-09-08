package hangman

import (
	"github.com/boltdb/bolt"
	"github.com/puffinframework/app"
)

const (
	appId = "Hangman"
)

type Hangman struct {
	db  *bolt.DB
	app *app.App
}

func NewApp(db *bolt.DB) *Hangman {

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

	return &Hangman{db: db, app: app}
}
