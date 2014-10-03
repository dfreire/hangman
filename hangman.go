package hangman

import (
	"github.com/dfreire/hangman/deltas"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	appId = "HangmanApp"
)

type HangmanApp struct {
	gormDB       gorm.DB
	deltaService deltas.DeltaService
}

func NewApp(deltaService deltas.DeltaService, gormDB gorm.DB) *HangmanApp {

	gormDB.CreateTable(Game{})

	return &HangmanApp{
		gormDB:       gormDB,
		deltaService: deltaService,
	}
}
