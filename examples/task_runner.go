package main

import (
	"github.com/jasocox/goworker"
	"log"
	"math/rand"
	"time"
)

var TASKS = 1000

type Task int

func (t Task) Do() error {
	log.Printf("Task #%d starting...", t)

	time.Sleep(time.Duration(int(time.Millisecond) * (rand.Intn(91) + 10)))

	log.Printf("Task #%d done.", t)

	return nil
}

func main() {
	worker := goworker.NewManager("Task Manager", 10)

	for i := 0; i < TASKS; i++ {
		worker.Exec(Task(i))
	}

	worker.Finish()
}
