package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func p1(filepath string) time.Duration {
	start := time.Now()
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var llist, rlist []int
	var lastpos int = 0

	for scanner.Scan() {
		// umm this is better ig Fields instead of -> // splittedline := strings.Split(scanner.Text(), " ")
		splittedline := strings.Fields(scanner.Text())
		if lastpos == 0 {
			lastpos = len(splittedline) - 1
		}

		lnum, err := strconv.Atoi(splittedline[0])
		if err != nil {
			log.Panic(err)
		}
		llist = append(llist, lnum)
		rnum, err := strconv.Atoi(splittedline[lastpos])
		if err != nil {
			log.Panic(err)
		}
		rlist = append(rlist, rnum)
	}

	slices.Sort(llist)
	slices.Sort(rlist)

	var diff int
	for i := 0; i < len(llist); i++ {
		diff += int(math.Abs(float64(llist[i] - rlist[i])))
	}
	log.Printf("Diff is %d", diff)

	return time.Since(start)
}

func p2(filepath string) time.Duration {
	start := time.Now()
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var llist []int
	var rlist []int
	var lastpos int = 0

	for scanner.Scan() {
		// umm this is better ig Fields instead of -> // splittedline := strings.Split(scanner.Text(), " ")
		splittedline := strings.Fields(scanner.Text())
		if lastpos == 0 {
			lastpos = len(splittedline) - 1
		}

		lnum, err := strconv.Atoi(splittedline[0])
		if err != nil {
			log.Panic(err)
		}
		llist = append(llist, lnum)
		rnum, err := strconv.Atoi(splittedline[lastpos])
		if err != nil {
			log.Panic(err)
		}
		rlist = append(rlist, rnum)
	}

	log.Println(llist)
	log.Println(rlist)
	lcount := make(map[int]int)
	lrOccur := make(map[int]int)
	for i := 0; i < len(llist); i++ {
		lcount[llist[i]]++
		lrOccur[rlist[i]]++
	}

	log.Println(lcount, lrOccur)

	var total int = 0

	for k, v := range lcount {
		log.Println("check", k, v)
		if lrOccur[k] > 0 {
			total += (k * lrOccur[k]) * v
		}
	}

	log.Printf("Total is %d", total)

	return time.Since(start)
}

func main() {
	t1 := p1("./day1/d1p1.txt")
	log.Printf("Time taken for p1 is %s", t1)

	t2 := p2("./day1/d1p1.txt")
	log.Printf("Time taken for p1 is %s", t2)
}
