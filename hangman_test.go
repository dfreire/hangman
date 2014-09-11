package hangman_test

import (
	"github.com/boltdb/bolt"
	. "github.com/dfreire/hangman"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	boltDB := openBoltDB()
	defer closeBoltDB(boltDB)

	gormDB := openGormDB()
	defer closeGormDB(gormDB)

	app := NewApp(boltDB, gormDB)

	evt, err := app.CreateGame(
		"TV",
		"Starring Steve Carell",
		"The Office",
		"http://en.wikipedia.org/wiki/The_Office",
		"dfreire",
	)
	assert.Nil(t, err)
	assert.NotNil(t, evt)

	game := evt.Data().(Game)
	assert.NotNil(t, game.Id)

	exists, err := app.ExistsGame(game.Id)
	assert.Nil(t, err)
	assert.True(t, exists)

	evt, err = app.RemoveGame(game.Id, game.AuthorId)
	assert.Nil(t, err)
	assert.NotNil(t, evt)

	exists, err = app.ExistsGame(game.Id)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func openBoltDB() *bolt.DB {
	db, err := bolt.Open("test-bolt.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	return db
}

func closeBoltDB(db *bolt.DB) {
	defer os.Remove(db.Path())
	db.Close()
}

func openGormDB() gorm.DB {
	db, err := gorm.Open("sqlite3", "./test-sql.db")
	if err != nil {
		panic(err)
	}
	db.SingularTable(true)
	return db
}

func closeGormDB(db gorm.DB) {
	defer os.Remove("./test-sql.db")
	db.Close()
}
