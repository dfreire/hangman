package hangman

import (
	"github.com/puffinframework/event"
)

const (
	CreatedGameEvent event.Type = "CreatedGameEvent"
	UpdatedGameEvent event.Type = "UpdatedGameEvent"
	RemovedGameEvent event.Type = "RemovedGameEvent"
)

func (self *Hangman) CreateGame(appId, theme, clue, answer, url, authorId string) (evt event.Event, err error) {
	evt = event.NewEvent(CreatedGameEvent, 1, Game{
		Theme:    theme,
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

func (self *Hangman) UpdateGame(appId, gameId, theme, clue, answer, url, authorId string) (evt event.Event, err error) {
	evt = event.NewEvent(UpdatedGameEvent, 1, Game{
		AppId:    appId,
		Id:       gameId,
		Theme:    theme,
		Clue:     clue,
		Answer:   answer,
		Url:      url,
		AuthorId: authorId,
	})
	return
}

func onUpdatedGame(evt event.Event) error {
	return nil
}

func (self *Hangman) OnUpdatedGame(evt event.Event) error {
	return onUpdatedGame(evt)
}

func (self *Hangman) RemoveGame(appId, gameId, authorId string) (evt event.Event, err error) {
	evt = event.NewEvent(RemovedGameEvent, 1, Game{
		AppId: appId,
		Id:    gameId,
	})
	return
}

func onRemovedGame(evt event.Event) error {
	return nil
}

func (self *Hangman) OnRemovedGame(evt event.Event) error {
	return onRemovedGame(evt)
}
