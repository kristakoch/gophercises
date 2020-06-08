package main

import (
	"flag"
)

var (
	directory = flag.String("quiz", "maths", "path to quiz files")                  // &q.name first arg
	version   = flag.String("version", "full", "designate short to limit to 5")     // &q.version
	order     = flag.String("order", "ordered", "enter rand to shuffle questions")  // &q.order
	seconds   = flag.Int("seconds", 99999999, "time limit for the quiz in seconds") // &q.seconds
)

func main() {
	flag.Parse()

	qn, qv, qo, qs := *directory, *version, *order, *seconds

	// Create the
	q := NewQuiz(qn, qv, qo, qs)

	// Deliver the
	DeliverQuiz(q)
}
