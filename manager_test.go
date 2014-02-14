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

func Test_ManagerDoesSeveralTasks(t *testing.T) {
	m := NewManager("Test Manaer", 5)

	vals := make([]*bool, 5)
	for i := range vals {
		vals[i] = new(bool)
	}

	for i := range vals {
		// No free i's around here....
		i := i
		m.Exec(NewTask(func() error {
			*vals[i] = true
			return nil
		}))
	}

	m.Finish()

	for _, val := range vals {
		if !*val {
			t.Error("Did not exec task")
			return
		}
	}
}
