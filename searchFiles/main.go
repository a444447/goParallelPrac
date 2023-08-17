package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

var (
	matches   []string
	waitgroup = sync.WaitGroup{}
	lock      = sync.Mutex{}
)

func searchFiles(root, filename string) {
	fmt.Printf("Searching in %s\n", root)
	files, _ := ioutil.ReadDir(root)
	for _, file := range files {
		if strings.Contains(file.Name(), filename) {
			lock.Lock()
			matches = append(matches, filepath.Join(root, file.Name()))
			lock.Unlock()
		}
		if file.IsDir() {
			waitgroup.Add(1)
			go searchFiles(filepath.Join(root, file.Name()), filename)
		}
	}
	waitgroup.Done()
}

func main() {
	waitgroup.Add(1)
	go searchFiles("E:\\GoLandProject\\goParallel", "note.md")
	waitgroup.Wait()
	for _, i := range matches {
		fmt.Printf("matched: %s\n", i)
	}
}
