package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func readCsv(filename string) (arr []Problem) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileReader := csv.NewReader(file)

	for {
		row, err := fileReader.Read()

		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		arr = append(arr, Problem{row[0], row[1]})
	}

	return
}

func startQuiz(questions []Problem, completeChan chan<- bool) {
	input := bufio.NewReader(os.Stdin)

	for _, q := range questions {
		log.Print("Question: ", q.question)
		ans, _ := input.ReadString('\n')
		log.Print("Answer: ", string(ans), q.answer)
	}

	completeChan <- true
}

func startTimer(limit int, completeChan chan<- bool) {
	if limit < 0 {
		return
	}

	<-time.After(time.Duration(limit) * time.Second)
	log.Print("Timer finished")
	completeChan <- true
}

func main() {
	// Define flags
	var (
		csvFilename = flag.String("csv", "problems.csv", "the location of the CSV file")
		timeLimit   = flag.Int("time-limit", 10, "the time limit for the user to answer all problems")
	)
	flag.Parse()

	// Read the csv given filename
	questionArr := readCsv(*csvFilename)
	completeChan := make(chan bool)

	go startTimer(*timeLimit, completeChan)
	go startQuiz(questionArr, completeChan)

	<-completeChan

	log.Print("Quiz complete")
}
