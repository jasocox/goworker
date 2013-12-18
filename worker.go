package goworker

import "errors"

type Worker interface {
  Name() string
  Messages() chan string
  Running() bool

  MustStart()
  Start() error
}

type worker struct {
  name string
  messages chan string
  running bool
  task func()
}

func NewWorker(name string, task func()) Worker {
  w := new(worker)
  w.name = name
  w.task = task
  w.messages = make(chan string)

  return w
}

func (w *worker) Name() string {
  return w.name
}

func (w *worker) Messages() chan string {
  return w.messages
}

func (w *worker) Running() bool {
  return w.running
}

func (w *worker) MustStart() {
  err := w.Start()
  if err != nil {
    panic(err.Error())
  }
}

func (w *worker) Start() (err error) {
  if w.running {
    err = errors.New(w.name + " is already running")
    return
  }

  w.running = true
  go func() {
    w.task()

    w.running = false
    w.messages <- "Done"
  }()

  return
}
