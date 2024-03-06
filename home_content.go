package main

import (
	"fmt"
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
	homeFiles.Append("Can't scan folder")
	homeFiles.Append("Go set a folder")
	homeWindow = w
	go func() {
		for {
			updateFiles()
			time.Sleep(30 * time.Second)
		}
	}()

	scroller := container.NewScroll(
		widget.NewListWithData(homeFiles,
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(di binding.DataItem, co fyne.CanvasObject) {
				co.(*widget.Label).Bind(di.(binding.String))
			},
		),
	)
	scroller.Direction = container.ScrollVerticalOnly

	return scroller
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
	canRead, err := storage.CanRead(uriFolder)
	check(err)
	if !canRead {
		fyne.CurrentApp().SendNotification(fyne.NewNotification("can't read folder", ""))
		return
	}
	canList, err := storage.CanList(uriFolder)
	check(err)

	if !canList {
		fmt.Println("Can't list")
		return
	}

	entries, err := storage.List(uriFolder)
	check(err)

	for _, entry := range entries {
		homeFiles.Append(entry.Name())
	}
}
