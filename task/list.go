package main

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your current tasks",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// // Open the database.
		// db, err := bolt.Open(DB_NAME, 0600, &bolt.Options{Timeout: 1 * time.Second})
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close()

		// // View the contents of the database.
		// err = db.Update(func(tx *bolt.Tx) error { // why am I using Update only to view tasks, not update them?

		// 	b := tx.Bucket([]byte(DB_BUCKET))

		// 	// List the tasks.
		// 	fmt.Println("You have the following tasks:")
		// 	count := 1
		// 	c := b.Cursor()
		// 	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		// 		fmt.Printf("%d. %s\n", count, k)
		// 		count++
		// 	}
		// 	if count == 1 {
		// 		fmt.Println("...no tasks for now!")
		// 	}

		// 	return nil
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
