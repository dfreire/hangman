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

type Delta struct {
	Header  DeltaHeader
	Records []DeltaRecord
}

type DeltaHeader struct {
	CommitNumber int
}

type DeltaRecord struct {
	Id            string
	Operation     DeltaOperation
	RecordType    string
	RecordVersion int
	Record        interface{}
}

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
		delta := Delta{
			Header: DeltaHeader{CommitNumber: 1},
			Records: []DeltaRecord{
				DeltaRecord{Operation: CREATE, RecordType: "Game", RecordVersion: 1, Record: game},
			},
		}
		headerBytes, _ := json.Marshal(delta.Header)
		recordsBytes, _ := json.Marshal(delta.Records)
		bucket.Put(headerBytes, recordsBytes)
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
		bucket := tx.Bucket([]byte(gamesBucketName))
		delta := Delta{
			Header: DeltaHeader{CommitNumber: 1},
			Records: []DeltaRecord{
				DeltaRecord{Operation: UPDATE, RecordType: "Game", RecordVersion: 1, Record: game},
			},
		}
		headerBytes, _ := json.Marshal(delta.Header)
		recordsBytes, _ := json.Marshal(delta.Records)
		bucket.Put(headerBytes, recordsBytes)
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
		bucket := tx.Bucket([]byte(gamesBucketName))
		delta := Delta{
			Header: DeltaHeader{CommitNumber: 1},
			Records: []DeltaRecord{
				DeltaRecord{Operation: REMOVE, RecordType: "Game", RecordVersion: 1, Record: game},
			},
		}
		headerBytes, _ := json.Marshal(delta.Header)
		recordsBytes, _ := json.Marshal(delta.Records)
		bucket.Put(headerBytes, recordsBytes)
		return self.OnDelta(delta)
	})
}

func (self *HangmanApp) OnDelta(delta Delta) error {
	for _, record := range delta.Records {
		if record.RecordType == "Game" {
			game := record.Record.(Game)
			switch record.Operation {
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
