package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

const (
	DB_NAME   = "tasks.db"
	DB_BUCKET = "TaskBucket"
)

func TestAdd(t *testing.T) {
	t.Run("normal add", func(t *testing.T) {
		taskName := "walk minnie"

		// Add the task.
		err, _ := runTask("add", taskName)
		if err != nil {
			t.Errorf("Error adding task '%s' to database.", taskName)
		}

		// Check to see if the task was successfully added.
		got, err := inDB(taskName)
		if err != nil {
			t.Errorf("Error checking whether task '%s' is in database.", taskName)
		}
		want := true

		if got != want {
			t.Errorf("Task '%s' not found in the database, should be found.", taskName)
		}

		// Remove the test task.
		err, _ = runTask("rm", taskName)
		if err != nil {
			t.Errorf("Error removing task '%s' from database.", taskName)
		}

	})
}

func TestDo(t *testing.T) {
	t.Run("normal do", func(t *testing.T) {
		taskName := "walk minnie"

		// Add the task.
		err, _ := runTask("add", taskName)
		if err != nil {
			t.Errorf("Error adding task '%s' to database.", taskName)
		}

		// Do the task.
		err, _ = runTask("do", taskName)
		if err != nil {
			t.Errorf("Error using do command for task '%s'.", taskName)
		}

		// Check to see if the item was successfully deleted.
		got, err := inDB(taskName)
		if err != nil {
			t.Errorf("Error checking whether task '%s' is in database.", taskName)
		}
		want := false

		if got != want {
			t.Errorf("Task '%s' found in the database, should be deleted.", taskName)
		}

	})
}

func TestList(t *testing.T) {
	// add 1 item
	// test the list command
	// delete the bucket

	// idea: flag for test mode which "runtask" can pass in and cmd
	//  will use the test bucket & delete it

	// Run list and get the results in a buffer.
	err, output := runTask("list", "")
	if err != nil {
		t.Errorf("Error adding task '%s' to database.", "")
	}
	fmt.Println("output of test list operation is:", output)

}

func TestRM(t *testing.T) {
	t.Run("normal remove", func(t *testing.T) {
		taskName := "walk minnie"

		// Add the task.
		err, _ := runTask("add", taskName)
		if err != nil {
			t.Errorf("Error adding task '%s' to database.", taskName)
		}

		// Remove the task.
		err, _ = runTask("rm", taskName)
		if err != nil {
			t.Errorf("Error removing task '%s' from database.", taskName)
		}

		// Check to see if the item was successfully deleted.
		got, err := inDB(taskName)
		if err != nil {
			t.Errorf("Error checking whether task '%s' is in database.", taskName)
		}
		want := false

		if got != want {
			t.Errorf("Task '%s' found in the database, should be deleted.", taskName)
		}

	})
}

func runTask(subCmd, taskName string) (error, string) {
	cmd := exec.Command("task", subCmd, taskName)
	var out, serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":" + fmt.Sprint(serr.String()))
		return err, ""
	}
	// fmt.Println(out.String()) // Prints command output if necessary.
	return nil, out.String()
}

func inDB(taskName string) (bool, error) {
	db, err := bolt.Open(DB_NAME, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte(DB_BUCKET))
	v := b.Get([]byte(taskName))

	// Commit the transaction and check for an error.
	if err := tx.Commit(); err != nil {
		return false, err
	}

	exists := false
	if v != nil {
		exists = true
	}

	return exists, nil
}
