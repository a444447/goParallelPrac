package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetter(url string, frequency *[26]int32, wg *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		index := strings.Index(allLetters, c)
		if index >= 0 {
			atomic.AddInt32(&frequency[index], 1)
		}
	}
	wg.Done()
}

func main() {
	var frequency [26]int32
	wg := sync.WaitGroup{}
	for i := 1000; i <= 1200; i++ {
		wg.Add(1)
		go countLetter("https://www.rfc-editor.org/rfc/rfc"+strconv.Itoa(i)+".txt", &frequency, &wg)
	}
	wg.Wait()
	for i, c := range allLetters {
		println(string(c), ":", frequency[i])
	}
}
