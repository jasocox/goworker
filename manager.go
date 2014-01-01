package goworker

import (
  "fmt"
  "log"
  "sync"
  "queue"
)

type Manager interface {
  Finish()
  Exec(Task)
}

type manager struct {
  name string

  available queue.Queue
  tasks queue.Queue
  wg sync.WaitGroup
}

func NewManager(name string, numWorkers int) Manager {
  log.Println("Making a new task manager: " + name)
  m := new(manager)
  m.name = name

  m.available = queue.New()
  m.tasks = queue.New()

  for i:=0; i<numWorkers; i++ {
    w := NewWorker(fmt.Sprintf("%s %d", name, i))
    m.available.Push(w)
  }

  return m
}

func (m *manager) Name() string {
  return m.name
}

func (m *manager) Finish() {
  log.Println(m.name + " is stopping")

  m.wg.Wait()
  for !m.available.IsEmpty() {
    m.available.Pop().(Worker).Stop()
  }

  return
}

func (m *manager) Exec(t Task) {
  var w Worker

  log.Println("Received a new task")
  if m.available.IsEmpty() {
    log.Println("No workers available, queueing")
    m.tasks.Push(t)
    return
  }

  w = m.available.Pop().(Worker)
  m.wg.Add(1)
  m.exec(w, t)
}

func (m *manager) exec(w Worker, t Task) {
  log.Println(m.name + ": Execing a task")
  w.Exec(t)
  go m.done(w)
}

func (m *manager) done(w Worker) {
  <-w.Messages()
  log.Println(m.name + ": Finished a task")

  if !m.tasks.IsEmpty() {
    var t Task
    t = m.tasks.Pop().(Task)
    m.exec(w, t)
  } else {
    m.wg.Done()
    m.available.Push(w)
  }
}
