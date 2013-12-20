package goworker

type Manager interface {
  Worker
}

type manager struct {
  w Worker

  workers []Worker
}

func NewManager(name string, numWorkers int, workerBuilder func() Worker) Manager {
  m := new(manager)
  m.workers = make([]Worker, numWorkers)

  m.w = NewWorker(name, func() {
    for _, worker := range m.workers {
      err := worker.Start()
      if err != nil {
        return
      }
    }

    m.waitForWorkers()
  })

  for i:=0; i<numWorkers; i++ {
    m.workers[i] = workerBuilder()
  }

  return m
}

func (m *manager) Name() string {
  return m.w.Name()
}

func (m *manager) Messages() chan string {
  return m.w.Messages()
}

func (m *manager) Running() bool {
  return m.w.Running()
}

func (m *manager) MustStart() {
  m.w.MustStart()
}

func (m *manager) Start() (err error) {
  return m.w.Start()
}

func (m *manager) waitForWorkers() {
  for _, worker := range m.workers {
    <-worker.Messages()
  }
}
