package goworker

import "log"

type Worker interface {
  Name() string
  Messages() <-chan string
  Stop()
  Exec(Task)
  Error() error
}

type worker struct {
  name string
  messages chan string
  commands chan string
  tasks chan Task
  err error
}

func NewWorker(name string) Worker {
  log.Println("Making a new worker: " + name)

  w := new(worker)
  w.name = name
  w.messages = make(chan string)
  w.commands = make(chan string)
  w.tasks = make(chan Task)

  go w.work()

  return w
}

func (w *worker) Name() string {
  return w.name
}

func (w *worker) Messages() <-chan string {
  return w.messages
}

func (w *worker) Stop() {
  log.Println("Stopping worker: " + w.name)

  w.commands <- "Quit"

  return
}

func (w *worker) Exec(t Task) {
  log.Println("Worker " + w.name + " is execing a task")
  w.tasks <- t
}

func (w *worker) Error() (err error) {
  err, w.err = w.err, nil

  return
}

func (w *worker) work() {
  for {
    select {
    case task := <-w.tasks:
      w.exec(task)
    case command := <-w.commands:
      if command == "Quit" {
        break;
      }
    }
  }

  log.Println("Worker " + w.name + " is stopping")
}

func (w *worker) exec(t Task) {
  err := t.Do()

  if err != nil {
    w.err = err
    w.messages <- "Error"
  } else {
    w.messages <- "Done"
  }
}
