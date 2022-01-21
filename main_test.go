package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	err := exec.Command("go", "build").Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	code := m.Run()
	os.Remove("find-by-extended-attribute")
	os.Exit(code)
}

func Test(t *testing.T) {
	t.Run("symlink", func(t *testing.T) {
		err := exec.Command("./fbea").Run()
		if err != nil {
			t.Fatal(err)
		}
	})
}
