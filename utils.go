package main

import (
	"github.com/charmbracelet/log"
)

func check(e error) {
	if e != nil {
		log.Error(e)
		panic(e)
		// fyne.CurrentApp().SendNotification(fyne.NewNotification("Error:", e.Error()))
	}
}
