package goworker

import (
	"errors"
	"testing"
)

type mockTask struct {
	ran bool
}

func (t *mockTask) Do() error {
	t.ran = true
	return nil
}

func (t *mockTask) DidRun() bool {
	return t.ran
}

func NewMockTask() (t *mockTask) {
	t = &mockTask{false}

	return
}

func Test_NewWorker(t *testing.T) {
	w := NewWorker("Test Worker")

	if w.Name() != "Test Worker" {
		t.Error("Didn't have the right name")
		return
	}
}

func Test_WorkerDoesTask(t *testing.T) {
	w := NewWorker("Test Worker")

	task := NewMockTask()
	w.Exec(task)

	if !task.DidRun() {
		t.Error("Did not exec the task")
	}
}

func Test_WorkerNotifiesWhenDone(t *testing.T) {
	w := NewWorker("Test Worker")

	w.Exec(NewMockTask())

	message := <-w.Messages()
	if message != DONE {
		t.Error("Did not get notified properly")
		return
	}
}

func Test_WorkerDoesMultipleTasks(t *testing.T) {
	w := NewWorker("Test Worker")

	task1 := NewMockTask()
	task2 := NewMockTask()
	task3 := NewMockTask()

	w.Exec(task1)
	message1 := <-w.Messages()
	w.Exec(task2)
	message2 := <-w.Messages()
	w.Exec(task3)
	message3 := <-w.Messages()

	if !task1.DidRun() {
		t.Error("Did not exec task 1")
		return
	}
	if !task2.DidRun() {
		t.Error("Did not exec task 2")
		return
	}
	if !task3.DidRun() {
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

	w.Exec(NewTask(func() error { return errors.New("Nope") }))

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

	task := NewMockTask()

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
