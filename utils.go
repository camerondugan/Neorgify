package main

import (
	"time"

	"github.com/charmbracelet/log"
)

type reminder struct {
	msg  string
	time time.Time
	file string
}

func sanitizeHash(b []byte) []byte {
	for i, v := range b {
		if v == '\n' { //new line
			b[i] = ' ' //not new line
		}
	}
	return b
}

func check(e error) {
	if e != nil {
		log.Error(e)
		panic(e)
		// fyne.CurrentApp().SendNotification(fyne.NewNotification("Error:", e.Error()))
	}
}
