package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io/fs"
	"os"
	"os/exec"
	"testing"
)

func Test(t *testing.T) {
	// build
	err := exec.Command("go", "build").Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// remove any existing test files
	// This ensures there are no incompatable files from executions of previous versions of the test squite.
	err = os.RemoveAll("test")
	if err != nil {
		panic(err)
	}

	// create test directory
	err = os.MkdirAll("test", fs.ModePerm)
	if err != nil {
		panic(err)
	}

	// create files with xattrs
	files := []struct {
		name       string
		xattrName  string
		xattrValue string
	}{
		// using random nouns because they are easy to recognise; makes debugging easier
		// without value
		{"gene", "pie", ""},
		{"artisan", "variety", ""},
		{"youth", "variety", ""},
		// with value
		{"activity", "mum", "climate"},
		{"event", "health", "security"},
		{"policy", "health", "security"},
	}
	for _, file := range files {
		// create file
		f, err := os.Create("test/" + file.name)
		if err != nil {
			panic(err)
		}
		// set xattr
		err = unix.Fsetxattr(int(f.Fd()), "user."+file.xattrName, []byte(file.xattrValue), 0)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	t.Cleanup(func() {
		// cleanup
		// remove executable
		err = os.Remove("find-by-extended-attribute")
		if err != nil {
			panic(err)
		}
	})

	t.Run("symlink", func(t *testing.T) {
		got, err := os.Readlink("fbea")
		if err != nil {
			t.Fatal(err)
		}
		want := "find-by-extended-attribute"
		if got != want {
			t.Fatalf("Want: %s, Got %s", want, got)
		}
	})

	t.Run("stdout", func(t *testing.T) {
		t.Run("justXattrName", func(t *testing.T) {
			t.Run("single", func(t *testing.T) {
				out, err := exec.Command("./find-by-extended-attribute", "user.pie").CombinedOutput()
				if err != nil {
					t.Log(string(out))
					t.Fatal(err)
				}
				got := string(out)
				want := "test/gene\n"
				if got != want {
					t.Fatalf("Want: %s, Got: %s", want, got)
				}
			})

			t.Run("multiple", func(t *testing.T) {
				out, err := exec.Command("./find-by-extended-attribute", "user.variety").CombinedOutput()
				if err != nil {
					t.Log(string(out))
					t.Fatal(err)
				}
				got := string(out)
				want := "test/artisan\ntest/youth\n"
				if got != want {
					t.Fatalf("Want: %s, Got: %s", want, got)
				}
			})
		})

		t.Run("xattrNameAndValue", func(t *testing.T) {
			t.Run("single", func(t *testing.T) {
				out, err := exec.Command("./find-by-extended-attribute", "user.mum", "climate").CombinedOutput()
				if err != nil {
					t.Log(string(out))
					t.Fatal(err)
				}
				got := string(out)
				want := "test/activity\n"
				if got != want {
					t.Fatalf("Want: %s, Got: %s", want, got)
				}
			})

			t.Run("multiple", func(t *testing.T) {
				out, err := exec.Command("./find-by-extended-attribute", "user.health", "security").CombinedOutput()
				if err != nil {
					t.Log(string(out))
					t.Fatal(err)
				}
				got := string(out)
				want := "test/policy\ntest/event\n"
				if got != want {
					t.Fatalf("Want: %s, Got: %s", want, got)
				}
			})
		})
	})
}
