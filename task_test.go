package goworker

import "testing"

func TestTaskCreation(t *testing.T) {
  NewTask(func() error {return nil})
}

func TestTaskRunning(t *testing.T) {
  var val *bool

  val = new(bool)

  task := NewTask(func() error {
    *val = true
    return nil
  })

  task.Do()

  if !*val {
    t.Error("Didn't exec the task")
  }
}

type MyTask struct {
  Val *bool
}

func (m *MyTask) Do() error {
  *m.Val = true

  return nil
}

func TestMyTask(t *testing.T) {
  m := MyTask{new(bool)}

  m.Do()

  if !*m.Val {
    t.Error("Did not run the task")
  }
}
