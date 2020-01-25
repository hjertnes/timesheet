// Package cmd sets up the command line interface
package cmd

import (
	"time"

	"github.com/hjertnes/timesheet/runner"
	"github.com/hjertnes/timesheet/utils"

	"github.com/spf13/cobra"
)

// RunFunc contains the runners used by Cobra
type RunFunc struct {
	r           runner.Runner
	ExcludedOpt bool
}

func (r *RunFunc) settingsList(cmd *cobra.Command, args []string) {
	r.r.SettingsList()
}
func (r *RunFunc) settingsSet(cmd *cobra.Command, args []string) {
	r.r.SettingsSet(args[0], args[1])
}
func (r *RunFunc) list(cmd *cobra.Command, args []string) {
	r.r.List()
}
func (r *RunFunc) off(cmd *cobra.Command, args []string) {
	date, err := utils.TimeFromDateString(args[0])
	utils.ErrorHandler(err)
	r.r.Off(date)
}
func (r *RunFunc) summaryDay(cmd *cobra.Command, args []string) {
	r.r.SummaryDay()
}
func (r *RunFunc) summaryYear(cmd *cobra.Command, args []string) {
	r.r.SummaryYear()
}
func (r *RunFunc) restore(cmd *cobra.Command, args []string) {
	r.r.Restore(args[0])
}
func (r *RunFunc) backup(cmd *cobra.Command, args []string) {
	r.r.Backup(args[0])
}
func (r *RunFunc) setup(cmd *cobra.Command, args []string) {
	r.r.Setup()
}
func (r *RunFunc) deleteOne(cmd *cobra.Command, args []string) {
	var date, err = utils.IntFromString(args[0])

	utils.ErrorHandler(err)
	r.r.Delete(date)
}
func (r *RunFunc) add(cmd *cobra.Command, args []string) {
	var err error

	var start time.Time

	var end time.Time

	start, err = utils.TimeFromDateStringAndTimeString(args[0], args[1])
	utils.ErrorHandler(err)

	end, err = utils.TimeFromDateStringAndTimeString(args[0], args[2])
	utils.ErrorHandler(err)

	r.r.Add(start, end, r.ExcludedOpt)
}

// New constructor
func New(run runner.Runner) *RunFunc {
	return &RunFunc{r: run}
}

type builder struct {
	run *RunFunc
}

func (b *builder) root() *cobra.Command {
	return &cobra.Command{
		Use:  "timesheet",
		Long: "A command line utility to keep track of worked hours",
	}
}

func (b *builder) settings() *cobra.Command {
	return &cobra.Command{
		Use:   "setting [sub-command]",
		Short: "settings",
		Long:  "manage settings",
	}
}

func (b *builder) settingsList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list settings",
		Long:  "command to list current settings",
		Args:  cobra.ExactArgs(0),
		Run:   b.run.settingsList,
	}
}

func (b *builder) settingsSet() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "set or update settings",
		Long:  "command to add or update settings",
		Args:  cobra.ExactArgs(2),
		Run:   b.run.settingsSet,
	}
}

func (b *builder) list() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "lists events",
		Long:  "lists all events in the database",
		Args:  cobra.ExactArgs(0),
		Run:   b.run.list,
	}
}

func (b *builder) add() *cobra.Command {
	return &cobra.Command{
		Use:   "add [date] [from] [to]",
		Short: "add event",
		Long: `logs work on a given date between two timestamps, 
deducts break for each date unless --excluded is used. Formats: yyyy-mm-dd, hh:mm`,
		Args: cobra.ExactArgs(3),
		Run:  b.run.add,
	}
}
func (b *builder) delete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "delete event",
		Long:  "remove event from database",
		Args:  cobra.ExactArgs(1),
		Run:   b.run.deleteOne,
	}
}

func (b *builder) off() *cobra.Command {
	return &cobra.Command{
		Use:   "off [date]",
		Short: "add day off",
		Long:  "Logs [date] as a day off, effiently deducting a day of working hours",
		Args:  cobra.ExactArgs(1),
		Run:   b.run.off,
	}
}

func (b *builder) setup() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "set timesheet up",
		Long:  "configure required settings. It will replace existing settings but not other data",
		Args:  cobra.ExactArgs(0),
		Run:   b.run.setup,
	}
}

func (b *builder) backup() *cobra.Command {
	return &cobra.Command{
		Use:   "backup [filename]",
		Short: "backup database to filename",
		Long:  "will write a json export of the database to filename",
		Args:  cobra.ExactArgs(1),
		Run:   b.run.backup,
	}
}

func (b *builder) restore() *cobra.Command {
	return &cobra.Command{
		Use:   "restore [filename]",
		Short: "restore database from export",
		Long:  "restores (and replaces) data in database from a previous export",
		Args:  cobra.ExactArgs(1),
		Run:   b.run.restore,
	}
}

func (b *builder) summaryYear() *cobra.Command {
	return &cobra.Command{
		Use:   "summary",
		Short: "show summary",
		Long:  "show a summary per year of how many hours are logged versus how many are expected",
		Args:  cobra.ExactArgs(0),
		Run:   b.run.summaryYear,
	}
}

func (b *builder) summaryDay() *cobra.Command {
	return &cobra.Command{
		Use:   "day",
		Short: "show hours logged per day",
		Long: `shows a list of dates and how many hours and minutes I worked 
in a format that makes it easy to copy paste into our time tracking stuff at work`,
		Args: cobra.ExactArgs(0),
		Run:  b.run.summaryDay,
	}
}

// Run builds and runs command
func Run(run *RunFunc, runner runner.Runner) {
	var b = &builder{
		run: run,
	}
	//Commands
	var rootCmd = b.root()

	var settingsCmd = b.settings()

	var settingsListCmd = b.settingsList()

	var settingsSetCmd = b.settingsSet()

	var listCmd = b.list()

	var offCmd = b.off()

	var addCmd = b.add()

	var deleteCmd = b.delete()

	var setupCmd = b.setup()

	var backupCmd = b.backup()

	var restoreCmd = b.restore()

	var summaryCmd = b.summaryYear()

	var summaryDayCmd = b.summaryDay()

	addCmd.Flags().BoolVarP(
		&run.ExcludedOpt,
		"excluded",
		"e",
		false,
		"will cause the days you use it on to not have break time deducted(e.g working extra hours during the weekend)",
	)

	settingsCmd.AddCommand(settingsListCmd)
	settingsCmd.AddCommand(settingsSetCmd)
	rootCmd.AddCommand(settingsCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(offCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	summaryCmd.AddCommand(summaryDayCmd)
	rootCmd.AddCommand(summaryCmd)
	_ = rootCmd.Execute()
}
