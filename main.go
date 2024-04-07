package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/markusmobius/go-dateparser"
)

type reminder struct {
	msg   string
	timer *time.Timer // controls when notification is sent
	time  time.Time   // time timer goes off
	file  string
}

var notesFolder string
var reminders []reminder
var parseConfig *dateparser.Configuration

func main() {
	notesFolder = getSettings("folder")
	reminders = []reminder{}
	parseConfig = &dateparser.Configuration{
		CurrentTime:         time.Now(),
		PreferredDateSource: dateparser.Future,
		Languages:           []string{"en"},
	}
	if notesFolder != "" {
		scanFolder(notesFolder)
	}
}

func setupReminder(r reminder) {
	if r.timer == nil && time.Now().Before(r.time) {
		whatsHappeningString := "created timer for " + r.msg + " at " + r.time.Format(
			"1:04PM on Mon 01-02-2006",
		)
		fmt.Println(whatsHappeningString)
		r.timer = time.AfterFunc(time.Until(r.time), func() {
			sendNtfy(r, "reminders")
		})
		// tell user the timer was created
		notifyTimerCreated := reminder{whatsHappeningString, nil, time.Now(), ""}
		sendNtfy(notifyTimerCreated, "timers")
	}
}

// https://docs.ntfy.sh/publish/?h=user#username-password
func sendNtfy(reminder reminder, board string) {
	server := getSettings("server")
	server += "/" + board
	req, err := http.NewRequest("POST", string(server), strings.NewReader(reminder.msg))
	check(err)
	secret := getSettings("login")
	req.Header.Set("Authorization", "Basic "+string(secret))
	_, err = http.DefaultClient.Do(req)
	check(err)
	log.Println("sent: " + reminder.msg)
}

func getSettings(file string) string {
	content, err := os.ReadFile(file)
	content = content[:len(content)-1]
	check(err)
	return string(content)
}
