package goworker

import "testing"

func Test_NewWorker(t *testing.T) {
  w := NewWorker("Test Worker")

  if w.Name() != "Test Worker" {
    t.Error("Didn't have the right name")
    return
  }
}

func Test_WorkerDoesTask(t *testing.T) {
  w := NewWorker("Test Worker")

  val := new(bool)
  w.Exec(NewTask(func() error {
    *val = true
    return nil
  }))

  if !*val {
    t.Error("Did not exec the task")
  }
}

func Test_WorkerDoesMultipleTasks(t *testing.T) {
  w := NewWorker("Test Worker")

  val1 := new(bool)
  val2 := new(bool)
  val3 := new(bool)

  task1 := NewTask(func() error {
    *val1 = true
    return nil
  })
  task2 := NewTask(func() error {
    *val2 = true
    return nil
  })
  task3 := NewTask(func() error {
    *val3 = true
    return nil
  })

  w.Exec(task1)
  <-w.Messages()
  w.Exec(task2)
  <-w.Messages()
  w.Exec(task3)
  <-w.Messages()

  if !*val1 {
    t.Error("Did not exec task 1")
  }
  if !*val2 {
    t.Error("Did not exec task 2")
  }
  if !*val3 {
    t.Error("Did not exec task 3")
  }
}
