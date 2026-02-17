package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deleteForce       bool
	deleteInteractive bool
	deleteList        string
	deleteFlagged     bool
)

var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a reminder",
	Example: `  rem delete abc12345
  rem rm abc12345 --force
  rem delete -i
  rem delete -i --list Work --flagged`,
	Args: func(cmd *cobra.Command, args []string) error {
		if deleteInteractive {
			return cobra.MaximumNArgs(0)(cmd, args)
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteInteractive {
			return runDeleteInteractive()
		}

		r, err := findReminderByID(args[0])
		if err != nil {
			return err
		}

		if !deleteForce {
			if isTTY() {
				confirmed, err := huhConfirm(fmt.Sprintf("Delete reminder '%s'?", r.Name))
				if err != nil {
					return err
				}
				if !confirmed {
					return nil
				}
			} else {
				return fmt.Errorf("use --force to delete non-interactively, or run in a terminal")
			}
		}

		if err := reminderSvc.DeleteReminder(r.ID); err != nil {
			return err
		}

		fmt.Printf("Deleted: %s\n", r.Name)
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation prompt")
	deleteCmd.Flags().BoolVarP(&deleteInteractive, "interactive", "i", false, "Select reminders interactively")
	deleteCmd.Flags().StringVarP(&deleteList, "list", "l", "", "Filter by list name")
	deleteCmd.Flags().BoolVar(&deleteFlagged, "flagged", false, "Filter to flagged reminders only")
	rootCmd.AddCommand(deleteCmd)
}

// runDeleteInteractive runs the interactive multi-select flow for deleting reminders.
func runDeleteInteractive() error {
	if err := requireInteractive(); err != nil {
		return err
	}

	reminders, err := reminderSvc.ListReminders(deleteFilter(deleteList, deleteFlagged))
	if err != nil {
		return err
	}

	selected, err := reminderMultiSelect("Select reminders to delete", reminders)
	if err != nil {
		return err
	}
	if selected == nil {
		return nil // cancelled
	}

	confirmed, err := huhConfirm(fmt.Sprintf("Delete %d reminder(s)?", len(selected)))
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	for _, id := range selected {
		if err := reminderSvc.DeleteReminder(id); err != nil {
			fmt.Printf("Error deleting %s: %v\n", shortIDStr(id), err)
			continue
		}
	}

	fmt.Printf("Deleted %d reminder(s)\n", len(selected))
	return nil
}
