package hangman

import (
	"github.com/boltdb/bolt"
	"github.com/puffinframework/app"
)

type Hangman struct {
	db  *bolt.DB
	app *app.App
}

func NewApp(db *bolt.DB) *Hangman {
	hangman := &Hangman{db: db, app: app.NewApp(db)}
	hangman.setup()
	return hangman
}

func (self *Hangman) setup() {
	exists, err := self.app.ExistsApp("Hangman")
	if err != nil {
		panic(err)
	}
	if !exists {
		self.app.CreateApp("Hangman")
	}
}
