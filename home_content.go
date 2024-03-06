package main

import (
	"errors"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func homeContent(window fyne.Window) fyne.CanvasObject {
	notesFolder := binding.NewString()
	notesFolder.Set("None Selected")
	onChoice := func(folderSelected fyne.ListableURI, err error) {
		if err != nil {
			errDialog := dialog.NewError(err, window)
			errDialog.Show()
		}
		if folderSelected != nil {
			log.Println("folder selected")
			log.Println(folderSelected.Name())
			validateFolder(folderSelected, window)
			notesFolder.Set(folderSelected.Path())
		} else {
			log.Println("no folder selected")
		}
	}
	folderDialog := dialog.NewFolderOpen(onChoice, window)
	onButton := func() {
		folderDialog.Show()
	}
	button := widget.NewButton("click me", onButton)
	folderLabel := widget.NewLabelWithData(notesFolder)
	return container.NewVBox(folderLabel, button)
}

// warns user of issues
func validateFolder(folderSelected fyne.ListableURI, window fyne.Window) {
	canRead, err := storage.CanRead(folderSelected)
	if err != nil {
		errDialog := dialog.NewError(err, window)
		errDialog.Show()
	}
	if !canRead {
		errDialog := dialog.NewError(errors.New("I can't read this folder! :<"), window)
		errDialog.Show()
	}

}

func notifyTest(a fyne.App) {
	notif := fyne.Notification{Title: "Notification", Content: "content"}
	for {
		a.SendNotification(&notif)
		time.Sleep(1 * time.Second)
	}
}
