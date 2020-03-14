// Package models for the new yaml based models
package models

import (
	"io/ioutil"
	"time"

	"github.com/hjertnes/timesheet/utils"
	"gopkg.in/yaml.v2"
)

// EventItem keeps track of a event with a start and end
type EventItem struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

// DayItem keeps track of a day and its ass events
type DayItem struct {
	Excluded bool        `yaml:"excluded,flow"`
	Events   []EventItem `yaml:"events"`
}

// Document is the root document structure
type Document struct {
	Configuration map[string]string             `yaml:"configuration,omitempty"`
	Items         map[string]map[string]DayItem `yaml:"items,inline"`
}

// Add an event
func (d *Document) Add(start time.Time, end time.Time, excluded bool, off bool) {
	year := start.Format("2006")
	day := start.Format("2006-01-02")

	yearItem, ok := d.Items[year]
	if !ok {
		yearItem = make(map[string]DayItem)
	}

	dayItem, ok := yearItem[day]
	if !ok {
		dayItem = DayItem{
			Excluded: false,
			Events:   make([]EventItem, 0),
		}
	}

	dayItem.Excluded = excluded

	if off {
		dayItem.Events = make([]EventItem, 0)
	} else {
		dayItem.Events = append(dayItem.Events, EventItem{
			Start: start.Format("15:04:05"),
			End:   end.Format("15:04:05"),
		})
	}

	yearItem[day] = dayItem
	d.Items[year] = yearItem
}

// Repository is the exposed interface
type Repository interface {
	Load() (*Document, error)
	Save(d *Document) error
}

type repository struct {
	filename string
}

func (r *repository) Load() (*Document, error) {
	f, err := utils.OpenOrCreate(r.filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var document Document

	err = yaml.Unmarshal(content, &document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (r *repository) Save(d *Document) error {
	f, err := utils.OpenOrCreate(r.filename)
	if err != nil {
		return err
	}

	defer f.Close()

	content, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	return err
}

func New(filename string) Repository {
	return &repository{
		filename,
	}
}
