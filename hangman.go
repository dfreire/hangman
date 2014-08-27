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
	a := app.NewApp(db)
	exists, err := a.ExistsApp("app1")
	if err != nil {
		panic(err)
	}
	if !exists {
		a.CreateApp("Hangman")
	}
	return &Hangman{db: db, app: a}
}
