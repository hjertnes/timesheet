// Package database wraps gorm.DB and some constuctor methods
package database

import (
	"fmt"
	"os"

	"github.com/hjertnes/timesheet/models"

	"github.com/jinzhu/gorm"
	//Sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Database type that contains gorm.DB object
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

// Open opens a regular database
func Open() *Database {
	var home = os.Getenv("HOME")

	var path = fmt.Sprintf("%s/.timesheet2.db", home)

	return &Database{Db: openDatabase(path)}
}

// OpenTest opens a database in current dir for testing purposes
func OpenTest() *Database {
	return &Database{Db: openDatabase("./data.db")}
}

// OpenInMemory opens a in memory database
func OpenInMemory() *Database {
	return &Database{Db: openDatabase(":memory:")}
}
