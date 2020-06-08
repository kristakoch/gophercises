package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Quiz type holds details about an instance of the quiz.
type Quiz struct {
	name      string
	questions [][]string
	version   string
	winImg    string
	loseImg   string
	quizLen   int
	order     string
	score     int
	seconds   int
}

// NewQuiz makes and returns a new Quiz
func NewQuiz(
	name string,
	version string,
	order string,
	seconds int,
) *Quiz {
	var q Quiz

	q.name = name
	q.version = version
	q.order = order
	q.seconds = seconds
	q.score = 0

	// Set win and lose images.
	setQuizAssets(&q)

	// Set quiz questions and length.
	setQuizQuestions(&q)

	return &q
}

// DeliverQuiz delivers a timed quiz from a Quiz.
func DeliverQuiz(q *Quiz) {
	welcomeMsg(q)
	deliverTimedQuiz(q)
	closingMsg(q)
}

// setQuizAssets attempts to return ascii images from win and lose text files in the quiz directory
func setQuizAssets(q *Quiz) error { // add error return
	winPath, losePath := fmt.Sprintf("%v/win.txt", q.name), fmt.Sprintf("%v/lose.txt", q.name)
	fWin, _ := ioutil.ReadFile(winPath)
	fLose, _ := ioutil.ReadFile(losePath)

	winStr, loseStr := string(fWin), string(fLose)
	q.winImg, q.loseImg = winStr, loseStr

	return nil
}

func setQuizQuestions(q *Quiz) {
	var questions [][]string
	filePath := fmt.Sprintf("%v/problems.csv", q.name)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fileStr := string(file)
	numQs := len(strings.Split(fileStr, "\n"))
	r := csv.NewReader(strings.NewReader(fileStr))
	for i := 0; i < numQs; i++ {
		record, err := r.Read()
		if err != nil {
			log.Fatal(err)
		}
		questions = append(questions, record)
	}
	if q.version == "short" && numQs > 5 {
		numQs = 5
	}

	q.questions = questions
	if q.order == "rand" {
		shuffleQuestions(q)
	}
	q.quizLen = numQs
	if q.version == "short" {
		q.questions = q.questions[:numQs]
	}
}

// deliverTimedQuiz delivers the quiz and times out after a given amount of seconds
func deliverTimedQuiz(q *Quiz) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nYou have %v seconds to complete this quiz. Press enter to begin...\n", q.seconds)
	reader.ReadString('\n')

	var err error
	done := make(chan bool, 1)
	go func() {
		err = deliverQuiz(q)
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(time.Duration(q.seconds) * time.Second):
	}
	if err != nil {
		return err
	}

	return nil
}

// deliverQuiz prompts the user for answers and returns a score of correct answers
func deliverQuiz(q *Quiz) error {
	reader := bufio.NewReader(os.Stdin)
	for num, ln := range q.questions {
		question := ln[0]
		answer := strings.ToLower(strings.Trim((ln[1]), "\n "))
		fmt.Printf("\n%v.%v: \n", num+1, question)
		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		text = strings.ToLower(strings.Trim(text, "\n "))
		matchedOne := false
		if strings.Contains(answer, "|") {
			matchedOne = answerMatch(text, answer)
		}
		if text == answer || matchedOne {
			// todo: display description
			color.Green("\n✓ correct...")
			fmt.Println("\n------")
			q.score++
		} else {
			// todo: display answer and description
			color.Red("\n⨉ nope\n")
			fmt.Println("\n------")
			fmt.Printf("correct answer: %v\n", answer)
		}
	}
	return nil
}

// welcomeMsg prints a message which includes the quiz name (name of the csv file)
func welcomeMsg(q *Quiz) {
	color.Blue("\n══════════☩═══✦═══☩══════════")
	fmt.Printf("\nWelcome to the %v quiz!\n\n", q.name)
	color.Blue("══════════☩═══✦═══☩══════════")
	fmt.Printf("\nThis quiz is %v questions.", q.quizLen)
	fmt.Println("\n\nLet's see if you know your facts...")
	fmt.Println("------")
}

// closingMsg prints a message which includes the score and, optionally, win/lose txt images
func closingMsg(q *Quiz) error { // add error return
	color.Blue("\n\n═════════════════════════════════")
	fmt.Printf("You scored %v/%v\n", q.score, q.quizLen)
	color.Blue("═════════════════════════════════")
	if q.score == q.quizLen {
		fmt.Print("\nCongratulations, you got em all!\n\n")
		fmt.Println(q.winImg)
	} else {
		fmt.Print("\nTry again...hehe\n\n")
		fmt.Println(q.loseImg)
	}
	return nil
}

// shuffleQuestions modifies quiz questions to be in random order
func shuffleQuestions(q *Quiz) error {
	origArr := make([][]string, len(q.questions))
	copy(origArr, q.questions)
	r, inc := rand.New(rand.NewSource(time.Now().Unix())), 0
	for _, i := range r.Perm(len(q.questions)) {
		q.questions[inc] = origArr[i]
		inc++
	}
	return nil
}

// answerMatch returns a boolean for whether  the answer matches one in a series of answers
func answerMatch(userResp string, answers string) bool {
	matchedOne := false
	allAnswers := strings.Split(answers, "|")
	answerBank := make(map[string]int)
	for i := 0; i < len(allAnswers); i++ {
		answerBank[allAnswers[i]] = 1
	}
	if answerBank[userResp] == 1 {
		matchedOne = true
	}
	return matchedOne
}
