package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var homeFiles binding.StringList
var homeWindow fyne.Window

func homeContent(w fyne.Window) fyne.CanvasObject {
	homeFiles = binding.NewStringList()

	setMessage()

	homeWindow = w
	go func() {
		for {
			updateFiles()
			time.Sleep(30 * time.Second)
		}
	}()

	itemList := widget.NewListWithData(homeFiles,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			co.(*widget.Label).Bind(di.(binding.String))
		},
	)
	scroller := container.NewScroll(itemList)
	scroller.Direction = container.ScrollVerticalOnly

	return scroller
}

// for when we can't scan folder
func setMessage() {
	homeFiles.Append("Can't scan folder")
	if fyne.CurrentApp().Preferences().StringWithFallback("NotesFolder", "") == "" {
		homeFiles.Append("No folder set")
	} else if fyne.CurrentDevice().IsMobile() {
		homeFiles.Append("Must give permission every time")
	} else {
		homeFiles.Append("Unknown Reason")
	}
}

func updateFiles() {
	if !pickedFolder && fyne.CurrentDevice().IsMobile() {
		return
	}
	homeFiles.Set([]string{})
	folder := fyne.CurrentApp().Preferences().StringWithFallback("NotesFolder", "")
	if folder == "" {
		time.Sleep(time.Second * 30)
		return
	}
	log.Println(folder)
	uriFolder, err := storage.ParseURI(folder)
	check(err)
	searchFolder(uriFolder)
}

func searchFolder(uriFolder fyne.URI) {
	canRead, err := storage.CanRead(uriFolder)
	check(err)
	if !canRead {
		fyne.CurrentApp().SendNotification(fyne.NewNotification("can't read folder", ""))
		return
	}

	isFolder, err := storage.CanList(uriFolder)
	check(err)
	if !isFolder {
		return
	}

	entries, err := storage.List(uriFolder)
	check(err)

	acceptableEndings := []string{".norg", ".md", ".txt", ".org"}
	filter := storage.NewExtensionFileFilter(acceptableEndings)
	for _, entry := range entries {
		if filter.Matches(entry) {
			homeFiles.Append(entry.String())
		}
		check(err)
		if !fyne.CurrentDevice().IsMobile() {
			searchFolder(entry)
		}
	}
}
