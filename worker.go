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
  debug bool
}

func NewWorker(name string) Worker {
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
  if w.debug {
    log.Printf("%s: stopping", w.name)
  }

  w.commands <- "Quit"

  return
}

func (w *worker) Exec(t Task) {
  if w.debug {
    log.Println("%s: received new task", w.name)
  }

  w.tasks <- t
}

func (w *worker) Error() (err error) {
  err, w.err = w.err, nil

  return
}

func (w *worker) SetDebug(d bool) {
  w.debug = d
}

func (w *worker) work() {
  for {
    select {
    case task := <-w.tasks:
      if w.debug {
        log.Printf("%s: execing a task", w.name)
      }
      w.exec(task)
    case command := <-w.commands:
      if command == "Quit" {
        break;
      }
    }
  }

  if w.debug {
    log.Printf("%s: is stopping", w.name)
  }
}

func (w *worker) exec(t Task) {
  err := t.Do()

  if err != nil {
    if w.debug {
      log.Printf("%s: error occured when execing task", w.name)
    }

    w.err = err
    w.messages <- "Error"
  } else {
    if w.debug {
      log.Printf("%s: done execing task")
    }

    w.messages <- "Done"
  }
}
