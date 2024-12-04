package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func readFile(filepath string, readRows chan<- []string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splittedline := strings.Fields(scanner.Text())
		readRows <- splittedline
	}
	close(readRows)
}

func convStrToInt(readRows <-chan []string, convertedRows chan<- []int, done chan struct{}) {
	for row := range readRows {
		pRow := make([]int, len(row))
		for i, v := range row {
			num, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("failed to convert: %v", err)
			}
			pRow[i] = num
		}
		convertedRows <- pRow
	}
	done <- struct{}{}
}

func checkSafeAfterConv_p1(processedRows <-chan []int, done chan struct{}, safeRowCount *atomic.Int64) {
	var increase, decrease, unsafe bool

	for row := range processedRows {
		increase, decrease, unsafe = false, false, false
		for i, _ := range row {
			if i == 0 {
				continue
			}
			if row[i-1] < row[i] {
				if decrease {
					// uh ig skip this shit to next row
					// log.Println("unsafe cause it inc while dec", row)
					unsafe = true
					break
				}
				if row[i]-row[i-1] > 3 {
					// now not safe, so skip to next row
					// log.Println("unsafe by val", row)
					unsafe = true
					break
				}
				increase = true
			} else if row[i-1] > row[i] {
				if increase {
					// log.Println("unsafe cause it dec while inc", row)
					unsafe = true
					break
				}
				if row[i-1]-row[i] > 3 {
					// log.Println("unsafe by val", row)
					unsafe = true
					break
				}
				decrease = true
			} else {
				// ig this is for equal
				unsafe = true
				break
			}
		}

		if !unsafe {
			// log.Println("Safe", row)
			safeRowCount.Add(1)
		}

	}
	done <- struct{}{}
}

func checkSafeAfterConv_p2(processedRows <-chan []int, done chan struct{}, safeRowCount *atomic.Int64) {
	for row := range processedRows {
		if isSequenceSafe(row) {
			safeRowCount.Add(1)
			continue
		}

		// um so brute force check all combinations if we find 1 comb to unsfe we break it
		// and we ig skip one one to check jhmmm
		for i := 0; i < len(row); i++ {
			newRow := make([]int, 0, len(row)-1)
			newRow = append(newRow, row[:i]...)
			newRow = append(newRow, row[i+1:]...)

			if isSequenceSafe(newRow) {
				safeRowCount.Add(1)
				break
			}
		}
	}
	done <- struct{}{}
}

func isSequenceSafe(row []int) bool {
	if len(row) < 2 {
		return true
	}

	isIncreasing := row[1] > row[0]

	for i := 1; i < len(row); i++ {
		diff := row[i] - row[i-1]

		if isIncreasing && diff <= 0 {
			return false
		}
		if !isIncreasing && diff >= 0 {
			return false
		}

		if abs(diff) < 1 || abs(diff) > 3 {
			return false
		}
	}

	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(filepath string, checkSafe func(processedRows <-chan []int, done chan struct{}, safeRowCount *atomic.Int64)) {
	start := time.Now()

	readRows := make(chan []string, 100)
	convertedRows := make(chan []int, 100)
	var safeRowCount atomic.Int64

	numRoutines := 3
	conversionDone := make(chan struct{}, numRoutines)
	processingDone := make(chan struct{}, numRoutines)

	go readFile(filepath, readRows)

	for i := 0; i < numRoutines; i++ {
		go convStrToInt(readRows, convertedRows, conversionDone)
	}

	for i := 0; i < numRoutines; i++ {
		go checkSafe(convertedRows, processingDone, &safeRowCount)
	}

	for i := 0; i < numRoutines; i++ {
		<-conversionDone
	}

	close(convertedRows)
	for i := 0; i < numRoutines; i++ {
		<-processingDone
	}

	println("Total safe: ", safeRowCount.Load())
	println("Total time: ", time.Since(start).String())
}

func main() {
	solve("./day2/d2p1.txt", checkSafeAfterConv_p1)
	solve("./day2/d2p1.txt", checkSafeAfterConv_p2)
}
