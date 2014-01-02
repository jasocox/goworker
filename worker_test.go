package goworker

import (
  "testing"
  "errors"
)

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

func Test_WorkerNotifiesWhenDone(t *testing.T) {
  w := NewWorker("Test Worker")

  w.Exec(NewTask(func() error {return nil}))

  message := <-w.Messages()
  if message != DONE {
    t.Error("Did not get notified properly")
    return
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
  message1 := <-w.Messages()
  w.Exec(task2)
  message2 := <-w.Messages()
  w.Exec(task3)
  message3 := <-w.Messages()

  if !*val1 {
    t.Error("Did not exec task 1")
    return
  }
  if !*val2 {
    t.Error("Did not exec task 2")
    return
  }
  if !*val3 {
    t.Error("Did not exec task 3")
    return
  }

  if message1 != DONE {
    t.Error("Did not get notified properly")
    return
  }
  if message2 != DONE {
    t.Error("Did not get notified properly")
    return
  }
  if message3 != DONE {
    t.Error("Did not get notified properly")
    return
  }
}

func Test_HandlesError(t *testing.T) {
  w := NewWorker("Test Worker")

  w.Exec(NewTask(func() error {return errors.New("Nope")}))

  message := <-w.Messages()

  if message != ERROR {
    t.Error("Did not respond with the correct status message")
    return
  }

  if w.Error().Error() != "Nope" {
    t.Error("Did not properly save the error")
    return
  }
}

func Test_StartingAndStopping(t *testing.T) {
  w := NewWorker("Test Worker")

  task := NewTask(func() error{return nil})

  err := w.Start()
  if err == nil {
    t.Error("Should have failed to start a running worker")
    return
  }

  err = w.Exec(task)
  if err != nil {
    t.Error("Shouldn't have failed to run this task")
    return
  }
  <-w.Messages()

  err = w.Stop()
  if err != nil {
    t.Error("Stopping should not have caused an error")
    return
  }

  err = w.Exec(task)
  if err == nil {
    t.Error("Execing using a stopped worker should have caused an error")
    return
  }

  err = w.Start()
  if err != nil {
    t.Error("Starting should not have caused an error")
    return
  }

  err = w.Exec(task)
  if err != nil {
    t.Error("Execing the task should not have caused an error")
    return
  }
  <-w.Messages()

  err = w.Start()
  if err == nil {
    t.Error("Starting the worker again should have caused an error")
    return
  }
}
