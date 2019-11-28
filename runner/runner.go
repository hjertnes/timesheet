package runner

import (
	"fmt"
	"time"
)

type IRunner interface {
	SettingsList()
	SettingsSet(key string, value string)
	List()
	Add(start time.Time, end time.Time, excluded bool)
	Off(date time.Time)
	Delete(id int)
	Setup()
	Backup(filename string)
	Restore(filename string)
	Import(filename string)
	SummaryYear()
	SummaryDay()
}

type Runner struct{}

func (Runner) SettingsList() {
	fmt.Println("Settings list")
}
func (Runner) SettingsSet(key string, value string) {
	fmt.Println("Settings set")
	fmt.Println(key)
	fmt.Println(value)
}
func (Runner) List() {
	fmt.Println("List")
}
func (Runner) Add(start time.Time, end time.Time, excluded bool) {
	fmt.Println("Add")
	fmt.Println(start)
	fmt.Println(end)
	fmt.Println(excluded)
}
func (Runner) Off(date time.Time) {
	fmt.Println("Off")
	fmt.Println(date)
}
func (Runner) Delete(id int) {
	fmt.Println("Delete")
	fmt.Println(id)

}
func (Runner) Setup() {
	fmt.Println("Setup")
}
func (Runner) Backup(filename string) {
	fmt.Println("Backup")
	fmt.Println(filename)
}
func (Runner) Restore(filename string) {
	fmt.Println("Restore")
	fmt.Println(filename)
}
func (Runner) Import(filename string) {
	fmt.Println("Import")
	fmt.Println(filename)
}
func (Runner) SummaryYear() {
	fmt.Println("Summary year")
}
func (Runner) SummaryDay() {
	fmt.Println("Summary day")
}

func NewRunner() IRunner {
	return &Runner{}
}
