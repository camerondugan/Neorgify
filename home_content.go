package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func homeContent(window fyne.Window) fyne.CanvasObject {
	onChoice := func(e fyne.ListableURI, err error) {
		if err != nil {
			panic(err)
		}
		log.Println("folder selected?")
	}
	folderDialog := dialog.NewFolderOpen(onChoice, window)
	onButton := func() {
		folderDialog.Show()
	}
	return widget.NewButton("click me", onButton)
}

func notifyTest(a fyne.App) {
	notif := fyne.Notification{Title: "Notification", Content: "content"}
	for {
		a.SendNotification(&notif)
		time.Sleep(1 * time.Second)
	}
}
