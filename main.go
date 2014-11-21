package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alinz/gonote/note"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	p := note.NewParser()
	if len(os.Args) != 2 {
		usage()
	}
	p.LoadFile(os.Args[1])
}
