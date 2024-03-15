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

func setupReminders() {
	for i := range reminders {
		if reminders[i].timer == nil && time.Now().Before(reminders[i].time) {
			fmt.Println("setup a timer for " + reminders[i].msg)
			reminders[i].timer = time.AfterFunc(reminders[i].time.Local().Sub(time.Now()), func() {
				sendNtfy(reminders[i])
			})
		}
	}
}

// https://docs.ntfy.sh/publish/?h=user#username-password
func sendNtfy(reminder reminder) {
	server := getSettings("server")
	req, err := http.NewRequest("POST", string(server), strings.NewReader(reminder.msg))
	check(err)
	secret := getSettings("login")
	req.Header.Set("Authorization", "Basic "+string(secret))
	http.DefaultClient.Do(req)
	log.Println("sent: " + reminder.msg)
}

func getSettings(file string) string {
	content, err := os.ReadFile(file)
	content = content[:len(content)-1]
	check(err)
	return string(content)
}
