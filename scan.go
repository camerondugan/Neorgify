package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/markusmobius/go-dateparser"
)

func scanFolder(folder string) {
	err := filepath.WalkDir(folder, walkDir)
	check(err)
	go setupReminders()
}

func walkDir(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if entry.IsDir() {
		return nil
	}
	acceptableEndings := []string{".norg", ".md", ".txt", ".org"}
	// detect if has acceptableEnding
	acceptable := false
	for _, ending := range acceptableEndings {
		if strings.HasSuffix(entry.Name(), ending) {
			acceptable = true
			break
		}
	}
	if acceptable {
		readFile(path)
	}
	return nil
}

func readFile(entry string) {
	filename := filepath.Base(entry)
	path := filepath.Dir(entry)

	fileBytes, err := os.ReadFile(entry)
	if !os.IsNotExist(err) {
		check(err)
	}

	h := sha256.Sum256(fileBytes)
	hash := sanitizeHash(h[:])
	hashPath := filepath.Join(path, ".neorgify_file_hashes")
	hashFile, err := os.Open(hashPath)
	defer hashFile.Close()
	if os.IsNotExist(err) { //if no hash or if we just started up: readTasks
		// create hash file for this folder
		f, err := os.Create(hashPath)
		check(err)
		defer f.Close()
		f.Write([]byte(filename + ":"))
		f.Write([]byte(hash))
		readTasksFromFile(fileBytes, entry)
		return
	}

	// hash file exists in this folder
	scanner := bufio.NewScanner(hashFile)
	var validLines []byte
	saveBuffer := bytes.NewBuffer(validLines)
	found := false
	for i := 0; scanner.Scan(); i++ {
		if len(scanner.Text()) == 0 {
			continue // remove empty lines
		}
		colonIndex := strings.IndexRune(scanner.Text(), ':')
		if colonIndex == -1 || len(scanner.Text()) == colonIndex {
			continue //remove invalid lines
		}
		storedHash := scanner.Text()[colonIndex+1:]
		if len(storedHash) != 32 { // invalid hash
			continue //remove invalid lines
		}
		hashedFileName := scanner.Text()[:colonIndex]
		if hashedFileName != filename {
			_, err := saveBuffer.Write(scanner.Bytes()) //save lines that don't match
			check(err)
			_, err = saveBuffer.Write([]byte("\n"))
			check(err)
			continue
		}
		if !slices.Equal([]byte(storedHash), hash) {
			// fmt.Printf("storedHash: %v\n", []byte(storedHash))
			// fmt.Printf("hash:       %v\n", hash)
			_, err := saveBuffer.Write([]byte(filename + ":")) //save new hash if doesn't match
			check(err)
			_, err = saveBuffer.Write([]byte(hash))
			check(err)
			_, err = saveBuffer.Write([]byte("\n"))
			check(err)
			readTasksFromFile(fileBytes, entry)
			continue
		}
		saveBuffer.Write(scanner.Bytes()) //save line if matches current hash
		saveBuffer.Write([]byte("\n"))
		// don't read tasks
		found = true
	}
	if !found {
		_, err := saveBuffer.Write([]byte(filename + ":")) //append new hash if doesn't exist
		check(err)
		_, err = saveBuffer.Write([]byte(hash))
		check(err)
		readTasksFromFile(fileBytes, entry)
	}
	if !slices.Equal(saveBuffer.Bytes(), fileBytes) {
		err := os.WriteFile(hashPath, saveBuffer.Bytes(), 0666)
		check(err)
	}
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
