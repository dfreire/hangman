package hangman_test

import (
	"github.com/boltdb/bolt"
	"github.com/dfreire/hangman"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	app := hangman.NewApp(db)

	evt, err := app.CreateGame(
		"Hangman",
		"TV",
		"Starring Steve Carell",
		"The Office",
		"http://en.wikipedia.org/wiki/The_Office",
		"dfreire",
	)
	assert.NotNil(t, evt)
	assert.Nil(t, err)
}

func openBoltDB() *bolt.DB {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	return db
}

func closeBoltDB(db *bolt.DB) {
	defer os.Remove(db.Path())
	db.Close()
}
