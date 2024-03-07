package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/markusmobius/go-dateparser"
)

var tasksWithDates binding.StringList
var homeWindow fyne.Window

func homeContent(w fyne.Window) fyne.CanvasObject {
	tasksWithDates = binding.NewStringList()

	setMessage()

	homeWindow = w
	go func() {
		for {
			updateFiles()
			time.Sleep(30 * time.Second)
		}
	}()

	itemList := widget.NewListWithData(tasksWithDates,
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
	tasksWithDates.Append("Can't scan folder")
	if fyne.CurrentApp().Preferences().StringWithFallback("NotesFolder", "") == "" {
		tasksWithDates.Append("No folder set")
	} else if fyne.CurrentDevice().IsMobile() {
		tasksWithDates.Append("Must give permission every time")
	} else {
		tasksWithDates.Append("Unknown Reason")
	}
}

func updateFiles() {
	if !pickedFolder && fyne.CurrentDevice().IsMobile() {
		return
	}
	tasksWithDates.Set([]string{})
	folder := fyne.CurrentApp().Preferences().StringWithFallback("NotesFolder", "")
	if folder == "" {
		time.Sleep(time.Second * 30)
		return
	}
	log.Println(folder)
	uriFolder, err := storage.ParseURI(folder)
	check(err)
	go searchFolder(uriFolder)
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
			readFile(entry)
		}
		check(err)
		if !fyne.CurrentDevice().IsMobile() {
			go searchFolder(entry)
		}
	}
}

func readFile(entry fyne.URI) {
	builder := strings.Builder{}
	taskPrefixes := []string{"- ( )"}
	parseConfig := &dateparser.Configuration{
		CurrentTime:         time.Now(),
		PreferredDayOfMonth: dateparser.First,
		PreferredDateSource: dateparser.Future,
	}
	// files
	closer, err := storage.Reader(entry)
	check(err)
	builder.Reset()
	buffer := make([]byte, 10000)
	var i = 1
	for i > 0 {
		i, _ = closer.Read(buffer)
		builder.Write(buffer)
	}
	for _, line := range strings.Split(builder.String(), "\n") {
		line := strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		for _, prefix := range taskPrefixes {
			if strings.HasPrefix(line, prefix) {
				line := line[len(prefix):]
				fmt.Printf("line: %v\n", line)
				_, dates, _ := dateparser.Search(parseConfig, line)
				if len(dates) > 0 {
					date := dates[len(dates)-1]
					tasksWithDates.Append(date.Date.Time.Format("Jan 2, 2006 at 3:04pm ") + line)
				}
			}
		}
	}
}
