package hangman

import (
	"encoding/json"
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

type Commit struct {
	Number int
	Deltas []Delta
}

type Delta struct {
	Id            string
	Operation     Operation
	RecordType    string
	RecordVersion int
	Record        interface{}
}

type Operation string

const (
	CREATE Operation = "CREATE"
	UPDATE Operation = "UPDATE"
	UPSERT Operation = "UPSERT"
	REMOVE Operation = "REMOVE"
)

type DeltaHandler interface {
	OnDelta(delta Delta)
}

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
		bucket := tx.Bucket([]byte(gamesBucketName))
		commit := Commit{
			Number: 1,
			Deltas: []Delta{
				Delta{Operation: CREATE, RecordType: "Game", RecordVersion: 1, Record: game},
			},
		}
		key, _ := json.Marshal(commit.Number)
		value, _ := json.Marshal(commit.Deltas)
		bucket.Put(key, value)
		return self.OnDeltas(commit.Deltas)
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
		bucket := tx.Bucket([]byte(gamesBucketName))
		commit := Commit{
			Number: 1,
			Deltas: []Delta{
				Delta{Operation: UPDATE, RecordType: "Game", RecordVersion: 1, Record: game},
			},
		}
		key, _ := json.Marshal(commit.Number)
		value, _ := json.Marshal(commit.Deltas)
		bucket.Put(key, value)
		return self.OnDeltas(commit.Deltas)
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
		bucket := tx.Bucket([]byte(gamesBucketName))
		commit := Commit{
			Number: 1,
			Deltas: []Delta{
				Delta{Operation: REMOVE, RecordType: "Game", RecordVersion: 1, Record: game},
			},
		}
		key, _ := json.Marshal(commit.Number)
		value, _ := json.Marshal(commit.Deltas)
		bucket.Put(key, value)
		return self.OnDeltas(commit.Deltas)
	})
}

func (self *HangmanApp) OnDeltas(deltas []Delta) error {
	for _, delta := range deltas {
		if delta.RecordType == "Game" {
			game := delta.Record.(Game)
			switch delta.Operation {
			case CREATE:
				self.gormDB.Create(game)
			case UPDATE:
				self.gormDB.Save(&game)
			case REMOVE:
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
