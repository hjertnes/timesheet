package cmd

import (
	"time"

	"github.com/hjertnes/timesheet/runner"
	"github.com/hjertnes/timesheet/utils"

	"github.com/spf13/cobra"
)

type RunFunc struct {
	r           runner.IRunner
	ExcludedOpt bool
}

func (r *RunFunc) SettingsList(cmd *cobra.Command, args []string) {
	r.r.SettingsList()
}
func (r *RunFunc) SettingsSet(cmd *cobra.Command, args []string) {
	r.r.SettingsSet(args[0], args[1])
}
func (r *RunFunc) List(cmd *cobra.Command, args []string) {
	r.r.List()
}
func (r *RunFunc) Off(cmd *cobra.Command, args []string) {
	date, err := utils.TimeFromDateString(args[0])
	utils.ErrorHandler(err)
	r.r.Off(date)
}
func (r *RunFunc) SummaryDay(cmd *cobra.Command, args []string) {
	r.r.SummaryDay()
}
func (r *RunFunc) SummaryYear(cmd *cobra.Command, args []string) {
	r.r.SummaryYear()
}
func (r *RunFunc) Restore(cmd *cobra.Command, args []string) {
	r.r.Restore(args[0])
}
func (r *RunFunc) Backup(cmd *cobra.Command, args []string) {
	r.r.Backup(args[0])
}
func (r *RunFunc) Setup(cmd *cobra.Command, args []string) {
	r.r.Setup()
}
func (r *RunFunc) Delete(cmd *cobra.Command, args []string) {
	var date, err = utils.IntFromString(args[0])
	utils.ErrorHandler(err)
	r.r.Delete(date)
}
func (r *RunFunc) Add(cmd *cobra.Command, args []string) {
	var err error
	var start time.Time
	var end time.Time

	start, err = utils.TimeFromDateStringAnTimeString(args[0], args[1])
	utils.ErrorHandler(err)

	end, err = utils.TimeFromDateStringAnTimeString(args[0], args[2])
	utils.ErrorHandler(err)

	r.r.Add(start, end, r.ExcludedOpt)
}
func NewRunFunc(run runner.IRunner) *RunFunc {
	return &RunFunc{r: run}
}

// Run test
func Run(run *RunFunc, runner runner.IRunner) {

	//Commands
	var rootCmd = &cobra.Command{
		Use:  "timesheet",
		Long: "A command line utility to keep track of worked hours",
	}
	var settingsCmd = &cobra.Command{
		Use:   "setting [sub-command]",
		Short: "settings",
		Long:  "manage settings",
	}
	var settingsListCmd = &cobra.Command{
		Use:   "list",
		Short: "list settings",
		Long:  "command to list current settings",
		Args:  cobra.ExactArgs(0),
		Run:   run.SettingsList,
	}
	var settingsSetCmd = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "set or update settings",
		Long:  "command to add or update settings",
		Args:  cobra.ExactArgs(2),
		Run:   run.SettingsSet,
	}
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "lists events",
		Long:  "lists all events in the database",
		Args:  cobra.ExactArgs(0),
		Run:   run.List,
	}
	var offCmd = &cobra.Command{
		Use:   "off [date]",
		Short: "add day off",
		Long:  "Logs [date] as a day off, effiently deducting a day of working hours",
		Args:  cobra.ExactArgs(1),
		Run:   run.Off,
	}
	//remember excluded flag
	var addCmd = &cobra.Command{
		Use:   "add [date] [from] [to]",
		Short: "add event",
		Long:  "logs work on a given date between two timestamps, deducts break for each date unless --exlcuded is used. Formats: yyyy-mm-dd, hh:mm",
		Args:  cobra.ExactArgs(3),
		Run:   run.Add,
	}
	var deleteCmd = &cobra.Command{
		Use:   "delete [id]",
		Short: "delete event",
		Long:  "remove event from database",
		Args:  cobra.ExactArgs(1),
		Run:   run.Delete,
	}
	var setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "set timesheet up",
		Long:  "configure required settings. It will replace existing settings but not other data",
		Args:  cobra.ExactArgs(0),
		Run:   run.Setup,
	}
	var backupCmd = &cobra.Command{
		Use:   "backup [filename]",
		Short: "backup database to filename",
		Long:  "will write a json export of the database to filename",
		Args:  cobra.ExactArgs(1),
		Run:   run.Backup,
	}
	var restoreCmd = &cobra.Command{
		Use:   "restore [filename]",
		Short: "restore database from export",
		Long:  "restores (and replaces) data in database from a previous export",
		Args:  cobra.ExactArgs(1),
		Run:   run.Restore,
	}
	//day option
	var summaryCmd = &cobra.Command{
		Use:   "summary",
		Short: "show summary",
		Long:  "show a summary per year of how many hours are logged versus how many are expected",
		Args:  cobra.ExactArgs(0),
		Run:   run.SummaryDay,
	}
	var summaryDayCmd = &cobra.Command{
		Use:   "day",
		Short: "show hours logged per day",
		Long:  "shows a list of dates and how many hours and minutes I worked in a format that makes it easy to copy paste into our time tracking stuff at work",
		Args:  cobra.ExactArgs(0),
		Run:   run.SummaryYear,
	}
	addCmd.Flags().BoolVarP(&run.ExcludedOpt, "excluded", "e", false, "will cause the days you use it on to not have break time deducted(e.g working extra hours during the weekend)")
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
	rootCmd.Execute()
}
