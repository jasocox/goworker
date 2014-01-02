package goworker

import (
  "log"
  "errors"
)

type Worker interface {
  Name() string
  SetDebug(bool)

  Start() error
  Stop() error

  Messages() <-chan int
  Exec(Task) error
  Error() error
}

type worker struct {
  name string

  messages chan int
  commands chan int
  tasks chan Task

  running bool
  err error
  debug bool
}

const (
  DONE = iota
  QUIT
  ERROR
)

func NewWorker(name string) Worker {
  w := new(worker)
  w.name = name
  w.messages = make(chan int)
  w.commands = make(chan int)
  w.tasks = make(chan Task)

  w.running = true
  go w.work()

  return w
}

func (w *worker) Name() string {
  return w.name
}

func (w *worker) Messages() <-chan int {
  return w.messages
}

func (w *worker) Start() error {
  if w.running {
    if w.debug {
      log.Printf("*s: already running", w.name)
    }

    return errors.New("Cannot start a running worker")
  }

  w.running = true
  go w.work()
  return nil
}

func (w *worker) Stop() error {
  if !w.running {
    if w.debug {
      log.Println("%s: is not running")
    }

    return errors.New("Cannot stop a worker that isn't working")
  }

  if w.debug {
    log.Printf("%s: stopping", w.name)
  }

  w.running = false
  w.commands <- QUIT

  return nil
}

func (w *worker) Exec(t Task) error {
  if !w.running {
    if w.debug {
      log.Println("%s: Cannot exec while not running")
    }
    return errors.New("Not running")
  }

  if w.debug {
    log.Println("%s: received new task", w.name)
  }
  w.tasks <- t

  return nil
}

func (w *worker) Error() (err error) {
  err, w.err = w.err, nil

  return
}

func (w *worker) SetDebug(d bool) {
  w.debug = d
}

func (w *worker) work() {
  if w.debug {
    log.Printf("%s: is starting")
  }

  for {
    select {
    case task := <-w.tasks:
      if w.debug {
        log.Printf("%s: execing a task", w.name)
      }
      w.exec(task)
    case command := <-w.commands:
      if command == QUIT {
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
    w.messages <- ERROR
  } else {
    if w.debug {
      log.Printf("%s: done execing task")
    }

    w.messages <- DONE
  }
}
