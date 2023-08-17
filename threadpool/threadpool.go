package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Point struct {
	x, y int
}

var (
	r               = regexp.MustCompile(`\((\d*),(\d*)\)`)
	numberOfThreads = 8
	waitGroup       = sync.WaitGroup{}
)

func findArea(inputChannel chan string) {
	for pointStr := range inputChannel {
		var points []Point
		for _, match := range r.FindAllStringSubmatch(pointStr, -1) {
			x, _ := strconv.Atoi(match[1])
			y, _ := strconv.Atoi(match[2])
			points = append(points, Point{x, y})
		}
		area := 0.0

		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)] // (i+1)%len 是因为最后一个点和第一个点相连
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		fmt.Println(math.Abs(area) / 2.0)
	}
	waitGroup.Done()
}

func main() {

	absPath, _ := filepath.Abs("threadpool/")
	dat, _ := ioutil.ReadFile(filepath.Join(absPath, "polygons.txt"))
	inputChannel := make(chan string, 100) //缓冲区大小
	for i := 0; i < numberOfThreads; i++ {
		go findArea(inputChannel)
	}
	waitGroup.Add(numberOfThreads)
	text := string(dat)
	for _, line := range strings.Split(text, "\n") {
		inputChannel <- line
	}
	close(inputChannel)
	waitGroup.Wait()
}
