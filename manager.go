package goworker

import (
	"fmt"
	"github.com/jasocox/figo"
	"log"
	"sync"
)

type Manager interface {
	Name() string
	Finish()
	Exec(Task)
	SetDebug(bool)
}

type manager struct {
	name string

	available figo.AsyncQueue
	tasks     figo.AsyncQueue
	wg        sync.WaitGroup

	debug bool
}

func NewManager(name string, numWorkers int) Manager {
	m := new(manager)
	m.name = name

	m.available = figo.NewAsync()
	m.tasks = figo.NewAsync()

	for i := 0; i < numWorkers; i++ {
		w := NewWorker(fmt.Sprintf("%s %d", name, i))
		m.available.Push(w)
	}

	return m
}

func (m *manager) SetDebug(d bool) {
	m.debug = d

	for elem := m.available.Front(); m.available.Next(elem) != nil; elem = m.available.Next(elem) {
		elem.Value.(Worker).SetDebug(d)
	}
}

func (m *manager) Name() string {
	return m.name
}

func (m *manager) Finish() {
	if m.debug {
		log.Println(m.name + " is stopping")
	}

	m.wg.Wait()
	for !m.available.IsEmpty() {
		m.available.Pop().(Worker).Stop()
	}

	return
}

func (m *manager) Exec(t Task) {
	var w Worker

	if m.debug {
		log.Println("Received a new task")
	}

	if m.available.IsEmpty() {
		if m.debug {
			log.Println("No workers available, queueing")
		}
		m.tasks.Push(t)
		return
	}

	w = m.available.Pop().(Worker)
	m.wg.Add(1)
	m.exec(w, t)
}

func (m *manager) exec(w Worker, t Task) {
	if m.debug {
		log.Println(m.name + ": Execing a task")
	}
	w.Exec(t)
	go m.done(w)
}

func (m *manager) done(w Worker) {
	<-w.Messages()
	if m.debug {
		log.Println(m.name + ": Finished a task")
	}

	if !m.tasks.IsEmpty() {
		var t Task
		t = m.tasks.Pop().(Task)
		m.exec(w, t)
	} else {
		m.wg.Done()
		m.available.Push(w)
	}
}
