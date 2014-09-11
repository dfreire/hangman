package hangman

import (
	"github.com/boltdb/bolt"
	_ "github.com/mattn/go-sqlite3"
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
	err = self.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		evt = event.NewEvent(CreatedGameEvent, 1, game)
		return onCreatedGame(b, evt)
	})
	return
}

func onCreatedGame(b *bolt.Bucket, evt event.Event) error {
	game := evt.Data().(Game)
	return b.Put([]byte(game.Id), []byte(game.AppId))
}

func (self *HangmanApp) OnCreatedGame(evt event.Event) error {
	return self.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		return onCreatedGame(b, evt)
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
	err = self.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		evt = event.NewEvent(RemovedGameEvent, 1, game)
		return onRemovedGame(b, evt)
	})
	return
}

func onRemovedGame(b *bolt.Bucket, evt event.Event) error {
	game := evt.Data().(Game)
	return b.Delete([]byte(game.Id))
}

func (self *HangmanApp) OnRemovedGame(evt event.Event) error {
	return self.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		return onRemovedGame(b, evt)
	})
}

func (self *HangmanApp) ExistsGame(gameId string) (exists bool, err error) {
	err = self.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		exists = existsGame(b, gameId)
		return nil
	})
	return
}

func existsGame(b *bolt.Bucket, gameId string) bool {
	return appId == string(b.Get([]byte(gameId)))
}
