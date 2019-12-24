// Package backupmodels models used to serialize / deserialize during export / import
package backupmodels

import "time"

// Document the enitre json document
type Document struct {
	Settings []Setting
	Events   []Event
}

// Setting model for setting
type Setting struct {
	Key   string
	Value string
}

// Event model for event
type Event struct {
	Start    time.Time
	End      time.Time
	Excluded bool
	Off      bool
}
