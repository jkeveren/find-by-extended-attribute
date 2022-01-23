package main

import (
	"bytes"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path/filepath"
)

func main() {
	l := log.New(os.Stdout, "", 0)                              // log
	le := log.New(os.Stderr, "find-by-extended-attribute: ", 0) // log error

	// handle bad input
	if len(os.Args) < 2 {
		le.Fatalf("Not enough arguments. Usage: find-by-extended-attribute <xattr name> [<xattr value>]")
	}

	// parse input
	xattrName := os.Args[1]
	var xattrValue []byte
	if len(os.Args) > 2 {
		xattrValue = []byte(os.Args[2])
	}

	// start recursion
	recurse(l, le, xattrName, xattrValue, ".")
}

func recurse(l, le *log.Logger, xattrName string, xattrValue []byte, p string) {
	// open file
	f, err := os.Open(p)
	if err != nil {
		le.Println(err)
		return
	}
	defer f.Close()

	// stat file
	info, err := f.Stat()
	if err != nil {
		le.Println(err)
		return
	}

	if info.IsDir() {
		// recurse deeper through directory
		// read directory
		names, err := f.Readdirnames(0)
		if err != nil {
			le.Println(err)
			return
		}

		// iterate through contents
		for _, name := range names {
			recurse(l, le, xattrName, xattrValue, filepath.Join(p, name))
		}
	} else {
		// check xattrs of file
		dest := make([]byte, len(xattrValue)) // len of 0 because most files will likely have no xattr value
		_, err := unix.Fgetxattr(int(f.Fd()), xattrName, dest)
		switch err {
		case nil:
			if bytes.Equal(xattrValue, dest) {
				l.Println(p)
			}
		case unix.ENODATA, unix.ERANGE:
		default:
			le.Println(err)
			return
		}
	}
}
