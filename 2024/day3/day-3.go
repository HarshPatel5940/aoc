package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// uh refer this -> https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_expressions/Cheatsheet

func solve_p1(r *regexp.Regexp, filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var total int = 0

	for scanner.Scan() {
		rawRes := r.FindAllStringSubmatch(scanner.Text(), -1)

		// log.Println(rawRes)
		for _, v := range rawRes {
			n1, err := strconv.Atoi(v[1])
			if err != nil {
				log.Fatal("skill issue", err)
			}

			n2, err := strconv.Atoi(v[2])
			if err != nil {
				log.Fatal("skill issue", err)
			}

			total += n1 * n2
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return total
}

func solve_p2(r *regexp.Regexp, filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var total int = 0

	var stop bool = false

	for scanner.Scan() {
		rawRes := r.FindAllStringSubmatch(scanner.Text(), -1)
		// log.Println("visiting", rawRes)

		for _, v := range rawRes {
			l := len(v)

			if strings.Contains(v[0], "do()") {
				// log.Println("doing")
				stop = false
			}

			if strings.Contains(v[0], "don't()") {
				// log.Println("not doing")
				stop = true
			}

			if stop {
				continue
			}

			if strings.Contains(v[0], "mul") {
				n1, err := strconv.Atoi(v[l-2])
				if err != nil {
					log.Fatal("skill issue", err)
				}

				n2, err := strconv.Atoi(v[l-1])
				if err != nil {
					log.Fatal("skill issue", err)
				}

				total += n1 * n2
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return total
}

func main() {
	p1 := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	log.Println("Ans for p1", solve_p1(p1, "./day3/d3p1.txt"))

	// never knew | is used for OR in regex thouht it was ||
	p2 := regexp.MustCompile(`(do\(\))|(don\'t\(\))|mul\((\d+),(\d+)\)`)
	log.Println("Ans for p2", solve_p2(p2, "./day3/d3p1.txt"))

}
