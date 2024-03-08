package main

import (
	"time"
)

var notesFolder string

func main() {
	notesFolder = "/home/cam/Notes"
	for {
		updateFiles()
	}
}

func updateFiles() {
	if notesFolder == "" {
		time.Sleep(time.Second * 30)
		return
	}
	searchFolder(notesFolder)
}
