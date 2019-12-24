// Package models contains database models
package models

import (
	"time"

	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Setting models for Settings
type Setting struct {
	ID    uint64 `gorm:"primary_key"`
	Key   string `gorm:"unique;not null"`
	Value string `gorm:"not null"`
}

// Event model for Events
type Event struct {
	ID       uint64    `gorm:"primary_key"`
	Start    time.Time `gorm:"not null"`
	End      time.Time
	Excluded bool `gorm:"not null"`
	Off      bool `gorm:"not null"`
}
