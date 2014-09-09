package hangman_test

import (
	"github.com/boltdb/bolt"
	. "github.com/dfreire/hangman"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	db := openBoltDB()
	defer closeBoltDB(db)

	app := NewApp(db)

	evt, err := app.CreateGame(
		"TV",
		"Starring Steve Carell",
		"The Office",
		"http://en.wikipedia.org/wiki/The_Office",
		"dfreire",
	)
	assert.Nil(t, err)
	assert.NotNil(t, evt)

	gameId := evt.Data().(Game).Id
	assert.NotNil(t, gameId)

	exists, err := app.ExistsGame(gameId)
	assert.Nil(t, err)
	assert.True(t, exists)
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
