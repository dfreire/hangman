package deltas

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

type DeltaService interface {
	Save(delta Delta, handler DeltaHandler) error
}

type BoltDeltaService struct {
	boltDB     *bolt.DB
	bucketName string
}

func NewBoltDeltaService(boltDB *bolt.DB, bucketName string) DeltaService {
	err := boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		panic(err)
	}
	return BoltDeltaService{
		boltDB:     boltDB,
		bucketName: bucketName,
	}
}

func (self BoltDeltaService) Save(delta Delta, handler DeltaHandler) error {
	return self.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(self.bucketName))

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
	})
}
