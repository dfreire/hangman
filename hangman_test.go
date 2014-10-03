package hangman_test

import (
	"github.com/boltdb/bolt"
	. "github.com/dfreire/hangman"
	"github.com/dfreire/hangman/deltas"
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

	deltaService := deltas.NewBoltDeltaService(boltDB, "HangmanDeltas")

	app := NewApp(deltaService, gormDB)

	evt, err := app.CreateGame(
		"TV",
		"Starring Steve Carell",
		"The Office",
		"http://en.wikipedia.org/wiki/The_Office",
		"dfreire",
	)
	assert.Nil(t, err)
	assert.NotNil(t, evt)

	game1 := evt.Data().(Game)
	assert.NotNil(t, game1.Id)

	exists, err := app.ExistsGame(game1.Id)
	assert.Nil(t, err)
	assert.True(t, exists)

	game2, err := app.GetGame(game1.Id)
	assert.Nil(t, err)
	assert.Equal(t, game1.Id, game2.Id)
	assert.Equal(t, game1.AppId, game2.AppId)
	assert.Equal(t, game1.Theme, game2.Theme)

	evt, err = app.RemoveGame(game1.Id)
	assert.Nil(t, err)
	assert.NotNil(t, evt)

	exists, err = app.ExistsGame(game1.Id)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestUpdate(t *testing.T) {
	boltDB := openBoltDB()
	defer closeBoltDB(boltDB)
	gormDB := openGormDB()
	defer closeGormDB(gormDB)

	deltaService := deltas.NewBoltDeltaService(boltDB, "HangmanDeltas")

	app := NewApp(deltaService, gormDB)

	evt1, err1 := app.CreateGame(
		"TV",
		"Starring Steve Carell",
		"The Office",
		"http://en.wikipedia.org/wiki/The_Office",
		"dfreire",
	)
	assert.Nil(t, err1)
	assert.NotNil(t, evt1)

	game1 := evt1.Data().(Game)
	assert.NotNil(t, game1.Id)

	evt2, err2 := app.UpdateGame(
		game1.Id,
		game1.Theme,
		"Starring Ricky Gervais",
		game1.Answer,
		game1.Url,
	)
	assert.Nil(t, err2)
	assert.NotNil(t, evt2)

	game2 := evt2.Data().(Game)
	assert.NotNil(t, game2.Id)

	assert.Equal(t, game1.Id, game2.Id)
	assert.Equal(t, game1.AppId, game2.AppId)
	assert.Equal(t, game1.Theme, game2.Theme)
	assert.NotEqual(t, game1.Clue, game2.Clue)

	game3, err := app.GetGame(game1.Id)
	assert.Nil(t, err)
	assert.Equal(t, game3.Id, game2.Id)
	assert.Equal(t, game3.AppId, game2.AppId)
	assert.Equal(t, game3.Theme, game2.Theme)
	assert.Equal(t, game3.Clue, game2.Clue)
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
	//db.LogMode(true)
	//db.SetLogger(gorm.Logger{revel.TRACE})
	//db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	return db
}

func closeGormDB(db gorm.DB) {
	defer os.Remove("./test-sql.db")
	db.Close()
}
