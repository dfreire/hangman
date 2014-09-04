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
	Id        string
	Front     string
	Back      string
	Url       string
	AuthorId  string
	Approved  bool
	FlagCount int
}

type Author struct {
	Id    string
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

func (self *Hangman) CreateCard(appId, front, back, url, authorId string) (evt event.Event, err error) {
	evt = event.NewEvent(CreatedCardEvent, 1, Card{
		Front:    front,
		Back:     back,
		Url:      url,
		AuthorId: authorId})
	return
}

func onCreatedCard(evt event.Event) error {
	return nil
}

func (self *Hangman) OnCreatedCard(evt event.Event) error {
	return onCreatedCard(evt)
}
