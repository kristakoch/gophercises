package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add something to your task list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		// Get the task.
		task := strings.Join(args, " ")

		// Put the task in the database.
		err := AddToDB(task)
		if err != nil {
			fmt.Print("Error when adding to DB:", err)
		}

		fmt.Printf("added '%s' to your task list\n", task)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

}
