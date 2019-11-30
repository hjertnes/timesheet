package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Setting struct {
	ID    uint64 `gorm:"primary_key"`
	Key   string `gorm:"unique;not null"`
	Value string `gorm:"not null"`
}

type Event struct {
	ID       uint64    `gorm:"primary_key"`
	Start    time.Time `gorm:"not null"`
	End      time.Time
	Excluded bool `gorm:"not null"`
	Off      bool `gorm:"not null"`
}
