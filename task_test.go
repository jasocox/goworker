package goworker

import "testing"

func TestTaskCreation(t *testing.T) {
  NewTask(func() error {return nil})
}
