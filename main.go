package main

import "github.com/alinz/gonote/note"

func main() {
	l := lexer.NewLexerWithFilename("test1", "examples/array-of-constant.note")
	l.Process()
}
