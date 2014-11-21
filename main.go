package main

import "github.com/alinz/gonote/note"

func main() {
	p := note.NewParser()
	p.LoadFile("examples/map.note")
}
