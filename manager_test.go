package goworker

import "testing"

func Test_New(t *testing.T) {
  m := NewManager("Test Manager", 0)

  if m.Name() != "Test Manager" {
    t.Error("Unexpected name")
    return
  }
}
