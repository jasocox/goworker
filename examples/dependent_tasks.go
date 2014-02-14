package main

import (
	"github.com/jasocox/goworker"
	"log"
	"math/rand"
	"time"
)

type task_one int

func (t task_one) Do() error {
  log.Printf("First task: %d starting...", t)

  time.Sleep(sleep_time())

  log.Printf("First task: %d done.", t)

  return nil
}

type task_two int

func (t task_two) Do() error {
  log.Printf("First task: %d starting...", t)

  time.Sleep(sleep_time())

  log.Printf("Second task: %d done.", t)

  return nil
}

var worker = goworker.NewManager("Dependent Tasks", 10)
const TASKS = 1000

func sleep_time() time.Duration {
  return time.Duration(int(time.Millisecond) * (rand.Intn(91) + 10))
}

func main() {
  for i := 0; i < TASKS; i++ {
    worker.Exec(task_one(i)).Then(task_two(i))
  }

  worker.Finish()
}
