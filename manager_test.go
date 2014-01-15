package goworker

import "testing"

func Test_New(t *testing.T) {
  m := NewManager("Test Manager", 0)

  if m.Name() != "Test Manager" {
    t.Error("Unexpected name")
    return
  }
}

func Test_ManagerDoesTask( t *testing.T) {
  m := NewManager("Test Manager", 1)

  val := new(bool)
  m.Exec(NewTask(func() error {
    *val = true
    return nil
  }))

  if !*val {
    t.Error("Did not run the task")
    return
  }
}
