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
	return &Hangman{db: db, app: app.NewApp(db)}
}
