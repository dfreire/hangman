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

type DeltaHeader struct {
	Id            string
	Operation     DeltaOperation
	RecordType    string
	RecordVersion int
}

type DeltaRecord interface{}

type DeltaOperation string

const (
	CREATE DeltaOperation = "CREATE"
	UPDATE DeltaOperation = "UPDATE"
	REMOVE DeltaOperation = "REMOVE"
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
		bucket := tx.Bucket([]byte(gamesBucketName))
		header := DeltaHeader{Operation: CREATE, RecordType: "Game", RecordVersion: 1}
		headerBytes, _ := json.Marshal(header)
		recordBytes, _ := json.Marshal(game)
		bucket.Put(headerBytes, recordBytes)
		return self.OnDelta(header, game)
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
		header := DeltaHeader{Operation: UPDATE, RecordType: "Game", RecordVersion: 1}
		headerBytes, _ := json.Marshal(header)
		recordBytes, _ := json.Marshal(game)
		bucket.Put(headerBytes, recordBytes)
		return self.OnDelta(header, game)
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
		header := DeltaHeader{Operation: REMOVE, RecordType: "Game", RecordVersion: 1}
		headerBytes, _ := json.Marshal(header)
		recordBytes, _ := json.Marshal(game)
		bucket.Put(headerBytes, recordBytes)
		return self.OnDelta(header, game)
	})
}

func (self *HangmanApp) OnDelta(header DeltaHeader, record DeltaRecord) error {
	game := record.(Game)
	switch header.Operation {
	case CREATE:
		self.gormDB.Create(game)
	case UPDATE:
		self.gormDB.Save(&game)
	case REMOVE:
		self.gormDB.Delete(&game)
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
