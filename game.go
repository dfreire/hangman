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

func (self *HangmanApp) CreateGame(theme, clue, answer, url, authorId string) (evt event.Event, err error) {
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

func (self *HangmanApp) OnCreatedGame(evt event.Event) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		return onCreatedGame(b, evt.Data().(Game))
	})
}

func (self *HangmanApp) UpdateGame(gameId, theme, clue, answer, url, authorId string) (evt event.Event, err error) {
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

func (self *HangmanApp) RemoveGame(gameId, authorId string) (evt event.Event, err error) {
	game := Game{
		AppId: appId,
		Id:    gameId,
	}
	err = self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		evt = event.NewEvent(RemovedGameEvent, 1, game)
		return onRemovedGame(b, game)
	})
	return
}

func onRemovedGame(b *bolt.Bucket, game Game) error {
	return b.Delete([]byte(game.Id))
}

func (self *HangmanApp) OnRemovedGame(evt event.Event) error {
	return self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		return onRemovedGame(b, evt.Data().(Game))
	})
}

func (self *HangmanApp) ExistsGame(gameId string) (exists bool, err error) {
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		exists = existsGame(b, gameId)
		return nil
	})
	return
}

func existsGame(b *bolt.Bucket, gameId string) bool {
	return appId == string(b.Get([]byte(gameId)))
}
