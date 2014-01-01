package goworker

import "testing"

func TestNewWorker(t *testing.T) {
  w := NewWorker("Test Worker")

  if w.Name() != "Test Worker" {
    t.Error("Didn't have the right name")
    return
  }
}
