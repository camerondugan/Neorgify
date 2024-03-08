package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/markusmobius/go-dateparser"
)

var notesFolder string
var startup bool
var reminders []reminder
var parseConfig *dateparser.Configuration

func main() {
	notesFolder = getSettings("folder")
	startup = true
	reminders = []reminder{}
	parseConfig = &dateparser.Configuration{
		CurrentTime:         time.Now(),
		PreferredDateSource: dateparser.Future,
		Languages:           []string{"en"},
	}
	// go func() {
	// 	for {
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()
	for {
		checkReminders()
		updateFiles()
		checkReminders()
		startup = false
		time.Sleep(10 * time.Second)
	}
}

func updateFiles() {
	if notesFolder != "" {
		searchFolder(notesFolder)
	}
}

func checkReminders() {
	curTime := time.Now()
	if !startup {
		for _, reminder := range reminders {
			if curTime.After(reminder.time) {
				//https://docs.ntfy.sh/publish/?h=user#username-password
				server := getSettings("server")
				req, err := http.NewRequest("POST", string(server), strings.NewReader(reminder.msg))
				check(err)
				secret := getSettings("login")
				req.Header.Set("Authorization", "Basic "+string(secret))
				http.DefaultClient.Do(req)
				fmt.Println("sent: " + reminder.msg)
			}
		}
	}
	// remove reminders we sent
	reminders = slices.DeleteFunc(reminders, func(r reminder) bool {
		return curTime.After(r.time)
	})
}

func getSettings(file string) string {
	content, err := os.ReadFile(file)
	content = content[:len(content)-1]
	check(err)
	return string(content)
}
