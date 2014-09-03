package hangman

import (
	"github.com/boltdb/bolt"
	"github.com/puffinframework/app"
	"github.com/puffinframework/event"
)

type Hangman struct {
	db  *bolt.DB
	app *app.App
}

const (
	CreatedCardEvent event.Type = "CreatedCardEvent"
)

type Card struct {
	Front     string
	Back      string
	Url       string
	Author    Author
	Approved  bool
	FlagCount int
}

type Author struct {
	email string
}

func NewHangman(db *bolt.DB) *Hangman {
	hangman := &Hangman{db: db, app: app.NewApp(db)}
	hangman.setup()
	return hangman
}

func (self *Hangman) setup() {
	exists, err := self.app.ExistsApp("hangman")
	if err != nil {
		panic(err)
	}
	if !exists {
		self.app.CreateApp("Hangman")
	}
}

func (self *Hangman) CreateCard(appId string) (evt event.Event, err error) {
	return nil, nil
}

func onCreatedCard(evt event.Event) error {
	return nil
}

func (self *Hangman) OnCreatedCard(evt event.Event) error {
	return onCreatedCard(evt)
}
