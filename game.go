package hangman

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/dfreire/hangman/deltas"
	"github.com/puffinframework/event"
	"github.com/satori/go.uuid"
)

const (
	deltasBucketName            = "HangmanDeltas"
	CreatedGameEvent event.Type = "CreatedGameEvent"
	UpdatedGameEvent event.Type = "UpdatedGameEvent"
	RemovedGameEvent event.Type = "RemovedGameEvent"
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
	err = self.OnCreatedGameEvent(evt)
	return
}

func (self *HangmanApp) OnCreatedGameEvent(evt event.Event) error {
	game := evt.Data().(Game)
	return self.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(deltasBucketName))
		delta := deltas.New([]deltas.Operation{
			deltas.Operation{Type: deltas.CREATE, Record: deltas.Record{Type: "Game", Version: 1, Value: game}},
		})
		key, _ := json.Marshal(delta.Id)
		value, _ := json.Marshal(delta.Operations)
		bucket.Put(key, value)
		return self.OnDelta(delta)
	})
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
	evt = event.NewEvent(UpdatedGameEvent, 1, game)
	err = self.OnUpdatedGameEvent(evt)
	return
}

func (self *HangmanApp) OnUpdatedGameEvent(evt event.Event) error {
	game := evt.Data().(Game)
	return self.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(deltasBucketName))
		delta := deltas.New([]deltas.Operation{
			deltas.Operation{Type: deltas.UPDATE, Record: deltas.Record{Type: "Game", Version: 1, Value: game}},
		})
		key, _ := json.Marshal(delta.Id)
		value, _ := json.Marshal(delta.Operations)
		bucket.Put(key, value)
		return self.OnDelta(delta)
	})
}

func (self *HangmanApp) RemoveGame(gameId string) (evt event.Event, err error) {
	game := Game{
		AppId: appId,
		Id:    gameId,
	}
	evt = event.NewEvent(RemovedGameEvent, 1, game)
	err = self.OnRemovedGameEvent(evt)
	return
}

func (self *HangmanApp) OnRemovedGameEvent(evt event.Event) error {
	game := evt.Data().(Game)
	return self.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(deltasBucketName))
		delta := deltas.New([]deltas.Operation{
			deltas.Operation{Type: deltas.REMOVE, Record: deltas.Record{Type: "Game", Version: 1, Value: game}},
		})
		key, _ := json.Marshal(delta.Id)
		value, _ := json.Marshal(delta.Operations)
		bucket.Put(key, value)
		return self.OnDelta(delta)
	})
}

func (self *HangmanApp) OnDelta(delta deltas.Delta) error {
	for _, operation := range delta.Operations {
		if operation.Record.Type == "Game" {
			game := operation.Record.Value.(Game)
			switch operation.Type {
			case deltas.CREATE:
				self.gormDB.Create(game)
			case deltas.UPDATE:
				self.gormDB.Save(&game)
			case deltas.REMOVE:
				self.gormDB.Delete(&game)
			}
		}
	}
	return nil
}

func (self *HangmanApp) GetGame(gameId string) (game Game, err error) {
	self.gormDB.Where("id = ?", gameId).First(&game)
	return
}

func (self *HangmanApp) ExistsGame(gameId string) (exists bool, err error) {
	game, err := self.GetGame(gameId)
	exists = game.Id == gameId
	return
}
