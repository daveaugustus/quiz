package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	csvFileName := flag.String("csv", "problems.csv", "A CSV file in form of Q&A")
	timeLimit := flag.Int("limit", 30, "The time limit for the quiz in seconds!")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		return
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		wg.Add(len(problems))
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
			wg.Done()
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				fmt.Println("Correct!")
				correct++
			} else {
				fmt.Println("Incorrect!")
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
	wg.Wait()
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
