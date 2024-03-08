package main

import "github.com/charmbracelet/log"

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
