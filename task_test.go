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

func NewErrorTask() Task {
	return NewTask(func() error {
		return errors.New("Nope")
	})
}

func TestTaskCreation(t *testing.T) {
	NewTask(func() error { return nil })
}

func TestTaskRunning(t *testing.T) {
	task := NewMockTask()

	if task.DidRun() {
		t.Error("Excuse me?")
		return
	}

	task.Do()
	if !task.DidRun() {
		t.Error("Didn't exec the task")
	}
}

func TestError(t *testing.T) {
	task := NewErrorTask()

	err := task.Do()

	if err == nil {
		t.Error("Did not get an error")
		return
	}

	if err.Error() != "Nope" {
		t.Error("Did not get the correct error")
		return
	}
}
