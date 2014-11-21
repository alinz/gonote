package main

import "github.com/alinz/gonote/note"

func main() {
	l, _ := note.NewLexerWithFilename("test1", "examples/map.note")
	l.Process()
}
