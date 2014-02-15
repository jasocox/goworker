package goworker

import "testing"

func Test_New(t *testing.T) {
	m := NewManager("Test Manager", 0)

	if m.Name() != "Test Manager" {
		t.Error("Unexpected name")
		return
	}
}

func Test_ManagerDoesTask(t *testing.T) {
	m := NewManager("Test Manager", 1)

	task := NewMockTask()
	m.Exec(task)

	if !task.DidRun() {
		t.Error("Did not run the task")
		return
	}
}

func Test_ManagerDoesSeveralTasks(t *testing.T) {
	m := NewManager("Test Manaer", 5)

	tasks := make([]*mockTask, 5)
	for i := range tasks {
		tasks[i] = NewMockTask()
	}

	for i := range tasks {
		m.Exec(tasks[i])
	}

	m.Finish()

	for _, task := range tasks {
		if !task.DidRun() {
			t.Error("Did not exec task")
			return
		}
	}
}

func Test_DependentTask(t *testing.T) {
	task1 := NewMockTask()
	task2 := NewMockTask()

	m := NewManager("Test Manager", 1)
	m.Exec(task1).Then(task2)
	m.Finish()

	if !task1.DidRun() {
		t.Error("Didn't run the first task")
		return
	}

	if !task2.DidRun() {
		t.Error("Didn't run the second task")
		return
	}
}
