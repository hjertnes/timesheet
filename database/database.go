package database

import (
	"fmt"
	"github.com/hjertnes/timesheet/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

type Database struct {
	Db *gorm.DB
}

func openDatabase(filename string) *gorm.DB {
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Setting{})
	db.AutoMigrate(&models.Event{})
	return db
}

func Open() *Database {
	var home = os.Getenv("HOME")
	var path = fmt.Sprintf("%s/.timesheet2.db", home)
	return &Database{Db: openDatabase(path)}

}

func OpenTest() *Database {
	return &Database{Db: openDatabase("./data.db")}
}

func OpenInMemory() *Database {
	return &Database{Db: openDatabase(":memory:")}
}
