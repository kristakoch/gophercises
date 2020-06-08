package main

import (
	"strings"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark something as done",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		task := strings.Join(args, " ")
		_ = task

		// // Open the database.
		// db, err := bolt.Open(DB_NAME, 0600, &bolt.Options{Timeout: 1 * time.Second})
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close()

		// // Delete the task from the database.
		// err = db.Update(func(tx *bolt.Tx) error {

		// 	b := tx.Bucket([]byte(DB_BUCKET))

		// 	// If the task doesn't exist, return an error.
		// 	if ok := b.Get([]byte(task)); ok == nil {
		// 		return fmt.Errorf("task '%s' not in task list", task)
		// 	}

		// 	if err := b.Delete([]byte(task)); err != nil {
		// 		return fmt.Errorf("remove from database error: %s", err)
		// 	}

		// 	fmt.Printf("hooray! you've done it: '%s'. removing from your task list\n", task)
		// 	return nil
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
