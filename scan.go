package main

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/markusmobius/go-dateparser"
)

func scanFolder(folder string) {
	watcher, err := fsnotify.NewWatcher()
	check(err)
	defer watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					//event.name is the file path
					log.Println("modified file at path:", event.Name)
					readIfAcceptable(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watcher err: ", err)
			}
		}
	}()

	// watch for every folder in the folder variable
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			log.Println(path)
			err = watcher.Add(path)
			check(err)
		} else { // files
			// read files on launch
			readIfAcceptable(path)
		}
		return nil
	})
	// never leave function
	<-make(chan struct{})
}

func readIfAcceptable(path string) {
	acceptableEndings := []string{".norg", ".md", ".txt", ".org"}
	// detect if has acceptableEnding
	acceptable := false
	for _, ending := range acceptableEndings {
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ending) {
			acceptable = true
			break
		}
	}
	if acceptable {
		readFile(path)
	}
}

func readFile(path string) {
	fileBytes, err := os.ReadFile(path)
	check(err)
	readTasksFromFile(fileBytes, path)
}

func readTasksFromFile(fileBytes []byte, path string) {
	deleteTasksFromMemory(path)
	taskPrefixes := []string{"( )", "[ ]"}
	for _, line := range strings.Split(string(fileBytes), "\n") {
		line := strings.Trim(line, " -")
		if len(line) == 0 {
			continue
		}
		for _, prefix := range taskPrefixes {
			if strings.HasPrefix(line, prefix) {
				line := line[len(prefix):]
				_, dates, _ := dateparser.Search(parseConfig, line)
				for _, date := range dates {
					message := strings.Replace(line, date.Text, "", 1)
					message = filepath.Base(path) + ": " + message
					reminder := reminder{msg: message, time: date.Date.Time.Local(), file: path}
					reminders = append(reminders, reminder)
					log.Println(date.Date.Time.Local().Format("Jan 2, 2006 at 3:04pm ") + line)
				}
			}
		}
	}
}

func deleteTasksFromMemory(file string) {
	shouldDelete := func(r reminder) bool {
		return r.file == file
	}

	// cancel timers
	for _, r := range reminders {
		if shouldDelete(r) {
			if r.timer != nil && !r.timer.Stop() {
				<-r.timer.C
			}
		}
	}

	reminders = slices.DeleteFunc(reminders, shouldDelete)
}
