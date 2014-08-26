package hangman

import (
	"github.com/boltdb/bolt"
)

type Aggregate struct {
	db *bolt.DB
}

func New(db *bolt.DB) *Aggregate {
	return &Aggregate{db: db}
}
