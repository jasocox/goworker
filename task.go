package goworker

type Task interface {
	Do() error
}

type task struct {
	t func() error
}

func NewTask(t func() error) Task {
	return &task{t}
}

func (t *task) Do() error {
	return t.t()
}
