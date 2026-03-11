package commands

import (
	"os"
	"time"

	"github.com/BRO3886/rem/internal/reminder"
	"github.com/BRO3886/rem/internal/ui"
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's and overdue reminders",
	Long:  `Show all incomplete reminders that are due today or overdue.`,
	Example: `  rem today
  rem today --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		endOfToday := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		completed := false
		filter := &reminder.ListFilter{
			Completed: &completed,
			DueBefore: &endOfToday,
		}

		reminders, err := reminderSvc.ListReminders(filter)
		if err != nil {
			return err
		}

		format := ui.ParseOutputFormat(outputFormat)
		ui.PrintReminders(os.Stdout, reminders, format)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
