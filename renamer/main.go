package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	pathpkg "path"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	// usage: nfc, err := RenameAll("marias+(.)*.txt", "marias birthday", "_")

	// Assign defualt bunk vals so files don't get blindly renamed in bulk.
	var re, np, s, d string
	flag.StringVar(&re, "regex", "defaultfn", "regex to match file names & extensions by")
	flag.StringVar(&np, "newprefix", "defaultfn", "name prefix to give to all files")
	flag.StringVar(&s, "separator", "_", "separates new prefix and number")
	flag.StringVar(&d, "directory", ".", "directory to recursively walk through")
	flag.Parse()

	nfc, err := RenameAll(re, np, s, d)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("files changed:", nfc)
}

// RenameAll renames files that match a regular expression
// in place and returns the number of files renamed.
func RenameAll(
	re string,
	newPrefix string,
	separator string,
	rDir string,
) (int, error) {
	var counter int
	err := filepath.Walk(rDir,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			fname := info.Name()
			matched, err := regexp.Match(re, []byte(fname))
			if err != nil {
				return err
			}

			if matched {
				// Build the new file path.
				dir, ext := pathpkg.Dir(path), pathpkg.Ext(path)
				numStr := strconv.Itoa(counter)
				newFP := "./" + dir + "/" + newPrefix + separator + numStr + ext

				fmt.Printf("renaming file %v to %v\n", info.Name(), newFP)

				// Rename here.
				err := os.Rename(path, newFP)
				if err != nil {
					return err
				}
				counter++
			}
			return nil
		})
	if err != nil {
		return 0, err
	}

	return counter, nil
}
