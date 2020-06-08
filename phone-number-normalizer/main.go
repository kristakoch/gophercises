package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbConn = "root@tcp(127.0.0.1)/meetupgroup"
)

func main() {
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the phone numbers from the db.
	q := "SELECT phone_number FROM user_info"
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		// Get the phone number value.
		var original string
		if err := rows.Scan(&original); err != nil {
			log.Fatal(err)
		}

		// Normalize the number.
		normalized, err := Normalize(original)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("normalizing data...", original, "--->", normalized)

		// Check the db to see if this number is already there.
		var col string
		q := "SELECT phone_number FROM user_info WHERE phone_number=?;"
		err = db.QueryRow(q, normalized).Scan(&col)

		// Check to see if a result was found or another
		// error occurred.
		var found bool
		if err != nil {
			if err == sql.ErrNoRows {
				found = false
			} else {
				log.Fatal(err)
			}
		} else {
			found = true
		}

		if found {
			// Already in db. Delete the original row.
			q = fmt.Sprintf(`
			DELETE FROM user_info
			WHERE phone_number='%v'`, original)

			fmt.Printf("normalized number %v already in db. deleting original row...\n\n", normalized)
		} else {
			// Not in db. Update the original row.
			q = fmt.Sprintf(`
			UPDATE user_info
			SET phone_number='%v'
			WHERE phone_number='%v'`, normalized, original)

			fmt.Printf("normalized number %v not in db. updating original row...\n\n", normalized)
		}

		// Run the delete or update query here.
		_, err = db.Exec(q)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Normalize takes in a phone number as a string and
// returns its normalized version.
func Normalize(pn string) (string, error) {
	es := ""
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	cutset := []string{" ", es, "(", es, ")", es, "{", es, "}", es, "-", es, "+", es, "/", es, ".", es}

	r := strings.NewReplacer(cutset...)
	res := r.Replace(pn)

	if strings.ContainsAny(strings.ToLower(res), alphabet) {
		return "", fmt.Errorf("invalid phone number: alpha chars present in '%s'", res)
	}

	return res, nil
}
