package hangman

import (
	"github.com/boltdb/bolt"
	"github.com/puffinframework/event"
	"github.com/satori/go.uuid"
)

const (
	gamesBucketName             = "HangmanGames"
	CreatedGameEvent event.Type = "CreatedGameEvent"
	UpdatedGameEvent event.Type = "UpdatedGameEvent"
	RemovedGameEvent event.Type = "RemovedGameEvent"
)

func (self *Hangman) CreateGame(appId, theme, clue, answer, url, authorId string) (evt event.Event, err error) {
	game := Game{
		Id:       uuid.NewV1().String(),
		AppId:    appId,
		Theme:    theme,
		Clue:     clue,
		Answer:   answer,
		Url:      url,
		AuthorId: authorId,
	}
	err = self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		evt = event.NewEvent(CreatedGameEvent, 1, game)
		return onCreatedGame(b, game)
	})
	return
}

func onCreatedGame(b *bolt.Bucket, game Game) error {
	return b.Put([]byte(game.Id), []byte(game.AppId))
}

func (self *Hangman) OnCreatedGame(evt event.Event) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		return onCreatedGame(b, evt.Data().(Game))
	})
}

func (self *Hangman) UpdateGame(appId, gameId, theme, clue, answer, url, authorId string) (evt event.Event, err error) {
	game := Game{
		Id:       gameId,
		AppId:    appId,
		Theme:    theme,
		Clue:     clue,
		Answer:   answer,
		Url:      url,
		AuthorId: authorId,
	}
	evt = event.NewEvent(UpdatedGameEvent, 1, game)
	return
}

func (self *Hangman) RemoveGame(appId, gameId, authorId string) (evt event.Event, err error) {
	evt = event.NewEvent(RemovedGameEvent, 1, Game{
		AppId: appId,
		Id:    gameId,
	})
	return evt, onRemovedGame(evt)
}

func onRemovedGame(evt event.Event) error {
	return nil
}

func (self *Hangman) OnRemovedGame(evt event.Event) error {
	return onRemovedGame(evt)
}
