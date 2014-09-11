package hangman_test

import (
    "database/sql"
	"github.com/boltdb/bolt"
	. "github.com/dfreire/hangman"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	boltDB := openBoltDB()
	defer closeBoltDB(boltDB)

    sqlDB := openSqlDB()
    defer closeSqlDB(sqlDB)

	app := NewApp(boltDB, sqlDB)

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

func openSqlDB() *sql.DB {
    db, err := sql.Open("sqlite3", "./test-sql.db")
	if err != nil {
		panic(err)
	}
	return db
}

func closeSqlDB(db *sql.DB) {
	defer os.Remove("./test-sql.db")
	db.Close()
}
