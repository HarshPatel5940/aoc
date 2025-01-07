package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type position struct {
	x int
	y int
}

type xmasSequence struct {
	positions [][2]int
}

func isInBounds(x, y, bound int) bool {
	return x >= 0 && x <= bound && y >= 0 && y <= bound
}

func fetchXPosition(rawArray [][]string, posiChannel chan position, wg *sync.WaitGroup) {
	defer wg.Done()
	for x, v := range rawArray {
		for y, val := range v {
			if strings.ToUpper(val) == "X" {
				posiChannel <- position{x: x, y: y}
			}
		}
	}
	close(posiChannel)
}

func checkAndPrintMatch(charArr [][]string, x, y int, dx, dy int, bound int) bool {
	if !isInBounds(x+3*dx, y+3*dy, bound) {
		return false
	}

	if charArr[x+dx][y+dy] == "M" &&
		charArr[x+2*dx][y+2*dy] == "A" &&
		charArr[x+3*dx][y+3*dy] == "S" {
		log.Printf("Found XMAS at (%d, %d) (%d, %d) (%d, %d) (%d, %d)", x, y, x+dx, y+dy, x+2*dx, y+2*dy, x+3*dx, y+3*dy)
		return true
	}
	return false
}

func fetchOccurances(posiChan <-chan position, bound int, charArr [][]string, count *atomic.Int32, wg *sync.WaitGroup) {
	defer wg.Done()

	directions := [][2]int{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	for pos := range posiChan {
		for _, dir := range directions {
			if checkAndPrintMatch(charArr, pos.x, pos.y, dir[0], dir[1], bound) {
				count.Add(1)
			}
		}
	}
}

func solve_p1(filePath string) {
	startTime := time.Now()
	file, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	var charArray [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(strings.ToUpper(scanner.Text()), "")
		charArray = append(charArray, line)
	}

	bound := len(charArray) - 1
	wg := sync.WaitGroup{}
	posiChan := make(chan position, 100)

	wg.Add(1)
	go fetchXPosition(charArray, posiChan, &wg)

	var xmasCount atomic.Int32
	wg.Add(2)

	go fetchOccurances(posiChan, bound, charArray, &xmasCount, &wg)
	go fetchOccurances(posiChan, bound, charArray, &xmasCount, &wg)

	wg.Wait()
	log.Printf("Found %d XMAS occurrences in %s", xmasCount.Load(), time.Since(startTime))
}

func main() {
	solve_p1("./day4/d4p1.txt")
}
