package main

// Usage:
// doug -> opens today in $EDITOR
// doug 2013-05-16
//      5-16
//      16-> opens that day in $EDITOR (uses current month / year if not provided)
// TODO: doug tomorrow, doug yesterday (flags?)
// TODO: custom environment variables: DOUG_EDITOR, DOUG_JOURNAL_HOME

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func setupDirpath(date *time.Time) (string, error) {
	dirpath := path.Join(
		os.Getenv("HOME"),
		"notes",
		"journal",
		date.Format("2006"),
		date.Format("01_January"))
	err := os.MkdirAll(dirpath, 0755)
	return dirpath, err
}

func getEditCmd() string {
	var result string
	if result = os.Getenv("DOUG_EDITOR"); result != "" {
		return result
	}
	if result = os.Getenv("EDITOR"); result != "" {
		return result
	}
	return "vi"
}

func editFile(filename string) {
	cmd := exec.Command(getEditCmd(), filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getJournalDate() time.Time {
	if len(os.Args) > 1 {
		result, err := parseTimeArg(os.Args[1])
		if err == nil {
			return result
		}
	}
	return time.Now()
}

// parses 1, 1-2, 1-2-2006
func parseTimeArg(timeArg string) (time.Time, error) {
	var year, month, day int
	switch strings.Count(timeArg, "-") {
	case 0:
		fmt.Sscanf(timeArg, "%2d", &day)
		now := time.Now()
		month = int(now.Month())
		year = now.Year()
	case 1:
		fmt.Sscanf(timeArg, "%2d-%2d", &month, &day)
		year = time.Now().Year()
	case 2:
		fmt.Sscanf(timeArg, "%2d-%2d-%4d", &month, &day, &year)
	default:
		now := time.Now()
		day = now.Day()
		month = int(now.Month())
		year = now.Year()
	}
	dateString := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	dateFormat := "2006-01-02"
	result, err := time.Parse(dateFormat, dateString)
	return result, err
}

func main() {
	journalDate := getJournalDate()
	dirpath, _ := setupDirpath(&journalDate)
	fileToEdit := path.Join(
		dirpath,
		journalDate.Format("02_Monday.org"))
	editFile(fileToEdit)
}
