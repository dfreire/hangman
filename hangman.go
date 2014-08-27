package hangman

import (
	"github.com/boltdb/bolt"
	"github.com/puffinframework/app"
)

type Hangman struct {
	db  *bolt.DB
	app *app.App
}

func NewHangman(db *bolt.DB) *Hangman {
	hangman := &Hangman{db: db, app: app.NewApp(db)}
	exists, err := hangman.app.ExistsApp("hangman")
	if err != nil {
		panic(err)
	}
	if !exists {
		hangman.app.CreateApp("Hangman")
	}
	return hangman
}
