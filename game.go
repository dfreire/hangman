package hangman

import (
	"github.com/puffinframework/event"
)

const (
	CreatedGameEvent event.Type = "CreatedGameEvent"
)

func (self *Hangman) CreateGame(appId, clue, answer, url, authorId string) (evt event.Event, err error) {
	evt = event.NewEvent(CreatedGameEvent, 1, Game{
		Clue:     clue,
		Answer:   answer,
		Url:      url,
		AuthorId: authorId,
	})
	return
}

func onCreatedGame(evt event.Event) error {
	return nil
}

func (self *Hangman) OnCreatedGame(evt event.Event) error {
	return onCreatedGame(evt)
}
