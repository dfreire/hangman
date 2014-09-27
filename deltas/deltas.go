package deltas

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
)

type Delta struct {
	Id         string
	Operations []Operation
}

type Operation struct {
	Type   OperationType
	Record Record
}

type OperationType string

const (
	CREATE OperationType = "CREATE"
	UPDATE OperationType = "UPDATE"
	UPSERT OperationType = "UPSERT"
	REMOVE OperationType = "REMOVE"
)

type Record struct {
	Type    string
	Id      string
	Version int
	Value   interface{}
}

type DeltaHandler func(delta Delta)

func New(operations []Operation) (delta Delta) {
	delta.Id = uuid.NewV1().String()
	delta.Operations = operations
	return
}

func withBoltDB(boltDB *bolt.DB, bucketName string) func(delta Delta, handler DeltaHandler) error {
	return func(delta Delta, handler DeltaHandler) error {
		return boltDB.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			return Save(bucket, delta, handler)
		})
	}
}

func Save(bucket *bolt.Bucket, delta Delta, handler DeltaHandler) error {
	key, err := json.Marshal(delta.Id)
	if err != nil {
		return err
	}

	value, err := json.Marshal(delta.Operations)
	if err != nil {
		return err
	}

	err = bucket.Put(key, value)
	if err != nil {
		return err
	}

	handler(delta)
	return nil
}
