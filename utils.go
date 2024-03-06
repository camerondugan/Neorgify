package main

import "fyne.io/fyne/v2"

func check(e error) {
	if e != nil {
		fyne.CurrentApp().SendNotification(fyne.NewNotification("Error:", e.Error()))
	}
}
