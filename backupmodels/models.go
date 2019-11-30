package backupmodels

import "time"

type Document struct {
	Settings []Setting
	Events []Event
}

type Setting struct {
	Key   string
	Value string
}

type Event struct {
	Start    time.Time
	End      time.Time
	Excluded bool
	Off      bool
}
