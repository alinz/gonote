package lexer

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

// Lexer holds the state of the scanner.
type Lexer struct {
	name    string     // used only for error reports.
	input   string     // the buffer reader to string content.
	start   int        // start position of this item.
	pos     int        // current position in the input.
	width   int        // width of last rune read from input.
	state   stateFn    // current state
	Tokens  chan Token // channel of scanned tokens.
	newLine bool       // once it reaches \n becomes true
}

type stateFn func(*Lexer) stateFn

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
	//I made width to zero so subsequent call to this method doesn't do any damage
	l.width = 0
}

// peek returns but does not consume
// the next rune in the input.
func (l *Lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

func (l *Lexer) currentSlice() string {
	return l.input[l.start:l.pos]
}

// next returns the next rune in the input.
// it only changes the pos and width for every call.
func (l *Lexer) next() rune {
	//if we reached the end of input it sets the width to zero
	//and return -1 which can be used to identified as end of string
	if l.pos >= len(l.input) {
		l.width = 0
		return -1
	}

	//tries to read the next rune. most of the case, width is 1.
	rune, size := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = size
	l.pos += l.width

	return rune
}

func (l *Lexer) current() rune {
	l.backup()
	return l.next()
}

// accept consumes the next rune
func (l *Lexer) accept(valid string) bool {
	val := strings.IndexRune(valid, l.next())
	if val != -1 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	//call next method until valid is not inside next item
	for strings.IndexRune(valid, l.next()) != -1 {
	}
	l.backup()
}

// acceptRunUtil consumes a run of runes until one of notValid set appears.
// so by calling current() you can find out which char in notValid string causes
// this to stop
func (l *Lexer) acceptRunUntil(notValid string) {
	for strings.IndexRune(notValid, l.next()) == -1 {
	}
	l.backup()
}

// emit emits a token based on internal pos variable
func (l *Lexer) emit(t TokenType) {
	l.Tokens <- Token{
		typ: t,
		val: l.input[l.start:l.pos],
	}
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.Tokens <- Token{tokenError, fmt.Sprintf(format, args...)}
	return nil
}

//NextToken gets a next token
func (l *Lexer) NextToken() Token {
	for {
		select {
		case tok := <-l.Tokens:
			return tok
		default:
			if l.state != nil {
				l.state = l.state(l)
			} else {
				return Token{typ: tokenEnd}
			}
		}
	}
}

func (l *Lexer) String() string {
	return fmt.Sprintf("-name: %s\n-input: \n%s-start: %d\n-pos: %d\n-width: %d\n", l.name, l.input, l.start, l.pos, l.width)
}

func (l *Lexer) indexSlice(lookFor rune) int {
	return strings.IndexRune(l.input[l.start:l.pos], lookFor)
}

//Process is a temporary call. I need to move it to Parser
func (l *Lexer) Process() {
	for {
		token := l.NextToken()

		fmt.Println(token)

		if token.typ == tokenEnd || token.typ == tokenError {
			break
		}

	}
}

//NewLexerWithString creates lexer based on string
func NewLexerWithString(name, input string) *Lexer {
	l := &Lexer{
		name:    name,
		input:   input,
		state:   LexDetect,
		Tokens:  make(chan Token, 2), // Two items sufficient.
		newLine: true,
	}

	return l
}

// NewLexerWithFilename create a lexer based on the file
func NewLexerWithFilename(name, input string) *Lexer {
	buf, err := ioutil.ReadFile(input)
	checkError(err)

	return NewLexerWithString(name, string(buf))
}
