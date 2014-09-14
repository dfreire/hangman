package hangman

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
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

type Delta struct {
	operation Operation
	deltaType string
	record    interface{}
}

type Operation string

const (
	CREATE Operation = "CREATE"
	UPDATE Operation = "UPDATE"
	REMOVE Operation = "REMOVE"
)

func (self *HangmanApp) CreateGame(theme, clue, answer, url, authorId string) (evt event.Event, err error) {
	id := uuid.NewV1()
	game := Game{
		Id:     id.String(),
		AppId:  appId,
		Theme:  theme,
		Clue:   clue,
		Answer: answer,
		Url:    url,
	}
	evt = event.NewEvent(CreatedGameEvent, 1, game)
	err = self.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(gamesBucketName))
		delta := Delta{operation: "CREATE", record: game}
		deltaBytes, _ := json.Marshal(delta)
		b.Put(id.Bytes(), deltaBytes)
		return self.OnDelta(delta)
	})
	return
}

func (self *HangmanApp) OnDelta(delta Delta) error {
	self.gormDB.Create(delta.record.(Game))
	return nil
}

func (self *HangmanApp) UpdateGame(gameId, theme, clue, answer, url string) (evt event.Event, err error) {
	game := Game{
		Id:     gameId,
		AppId:  appId,
		Theme:  theme,
		Clue:   clue,
		Answer: answer,
		Url:    url,
	}
	err = self.boltDB.Update(func(tx *bolt.Tx) error {
		evt = event.NewEvent(UpdatedGameEvent, 1, game)
		return onUpdatedGame(self.gormDB, evt)
	})
	return
}

func onUpdatedGame(gormDB gorm.DB, evt event.Event) error {
	game := evt.Data().(Game)
	gormDB.Save(&game)
	return nil
}

func (self *HangmanApp) OnUpdatedGame(evt event.Event) error {
	return onUpdatedGame(self.gormDB, evt)
}

func (self *HangmanApp) RemoveGame(gameId string) (evt event.Event, err error) {
	game := Game{
		AppId: appId,
		Id:    gameId,
	}
	evt = event.NewEvent(RemovedGameEvent, 1, game)
	err = onRemovedGame(self.gormDB, evt)
	return
}

func onRemovedGame(gormDB gorm.DB, evt event.Event) error {
	game := evt.Data().(Game)
	gormDB.Delete(&game)
	return nil
}

func (self *HangmanApp) OnRemovedGame(evt event.Event) error {
	return onRemovedGame(self.gormDB, evt)
}

func (self *HangmanApp) ExistsGame(gameId string) (exists bool, err error) {
	game := Game{}
	self.gormDB.Where("id = ?", gameId).First(&game)
	return game.Id == gameId, nil
}

func existsGame(b *bolt.Bucket, gameId string) bool {
	return appId == string(b.Get([]byte(gameId)))
}

func (self *HangmanApp) GetGame(gameId string) (game Game, err error) {
	self.gormDB.Where("id = ?", gameId).First(&game)
	return
}
