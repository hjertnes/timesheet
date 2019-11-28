package cmd

import (
	"time"

	"github.com/hjertnes/timesheet/runner"
	"github.com/hjertnes/timesheet/utils"

	"github.com/spf13/cobra"
)

// Run test
func Run(runner runner.IRunner) {
	var excludedOpt bool
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
		Run: func(cmd *cobra.Command, args []string) {
			runner.SettingsList()
		},
	}
	var settingsSetCmd = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "set or update settings",
		Long:  "command to add or update settings",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			runner.SettingsSet(args[0], args[1])
		},
	}
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "lists events",
		Long:  "lists all events in the database",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runner.List()
		},
	}
	var offCmd = &cobra.Command{
		Use:   "off [date]",
		Short: "add day off",
		Long:  "Logs [date] as a day off, effiently deducting a day of working hours",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			date, err := utils.TimeFromDateString(args[0])
			utils.ErrorHandler(err)
			runner.Off(date)
		},
	}
	//remember excluded flag
	var addCmd = &cobra.Command{
		Use:   "add [date] [from] [to]",
		Short: "add event",
		Long:  "logs work on a given date between two timestamps, deducts break for each date unless --exlcuded is used. Formats: yyyy-mm-dd, hh:mm",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var start time.Time
			var end time.Time

			start, err = utils.TimeFromDateStringAnTimeString(args[0], args[1])
			utils.ErrorHandler(err)

			end, err = utils.TimeFromDateStringAnTimeString(args[0], args[2])
			utils.ErrorHandler(err)

			runner.Add(start, end, excludedOpt)
		},
	}
	var deleteCmd = &cobra.Command{
		Use:   "delete [id]",
		Short: "delete event",
		Long:  "remove event from database",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var date, err = utils.IntFromString(args[0])
			utils.ErrorHandler(err)
			runner.Delete(date)
		},
	}
	var setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "set timesheet up",
		Long:  "configure required settings. It will replace existing settings but not other data",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runner.Setup()
		},
	}
	var backupCmd = &cobra.Command{
		Use:   "backup [filename]",
		Short: "backup database to filename",
		Long:  "will write a json export of the database to filename",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runner.Backup(args[0])
		},
	}
	var importCmd = &cobra.Command{
		Use:   "import [filename]",
		Short: "import data",
		Long:  "import data from my previous system",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runner.Import(args[0])
		},
	}
	var restoreCmd = &cobra.Command{
		Use:   "restore [filename]",
		Short: "restore database from export",
		Long:  "restores (and replaces) data in database from a previous export",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runner.Restore(args[0])
		},
	}
	//day option
	var summaryCmd = &cobra.Command{
		Use:   "summary",
		Short: "show summary",
		Long:  "show a summary per year of how many hours are logged versus how many are expected",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runner.SummaryYear()
		},
	}
	var summaryDayCmd = &cobra.Command{
		Use:   "day",
		Short: "show hours logged per day",
		Long:  "shows a list of dates and how many hours and minutes I worked in a format that makes it easy to copy paste into our time tracking stuff at work",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runner.SummaryDay()
		},
	}
	addCmd.Flags().BoolVarP(&excludedOpt, "excluded", "e", false, "will cause the days you use it on to not have break time deducted(e.g working extra hours during the weekend)")
	settingsCmd.AddCommand(settingsListCmd)
	settingsCmd.AddCommand(settingsSetCmd)
	rootCmd.AddCommand(settingsCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(offCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(restoreCmd)
	summaryCmd.AddCommand(summaryDayCmd)
	rootCmd.AddCommand(summaryCmd)
	rootCmd.Execute()
}
