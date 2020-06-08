/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"strings"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a task",
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

		// 	fmt.Printf("removed '%s' from your task list\n", task)
		// 	return nil
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }

	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
