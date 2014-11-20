package lexer

import "strings"

const (
	indentMeta = "\t"
	arrayMeta  = "- "
	spaceMeta  = ' '
	nullMeta   = "null"
	trueMeta   = "true"
	falseMeta  = "false"
	configMeta = "@"
	returnMeta = '\n'
	eof        = -1
)

func LexArrayInside(l *Lexer) stateFn {
	return nil
}

func LexArrayStart(l *Lexer) stateFn {
	l.pos += len(arrayMeta)
	l.emit(tokenArrayStart)
	return LexDetect
}

func LexConfig(l *Lexer) stateFn {
	l.pos += len(configMeta)
	l.start = l.pos

	l.acceptRunUntil("\n")

	l.emit(tokenConfig)

	return LexDetect
}

func LexEnter(l *Lexer) stateFn {
	l.emit(tokenEnter)
	l.newLine = true
	return LexDetect
}

func LexEnd(l *Lexer) stateFn {
	l.width = 1
	l.emit(tokenEnd)
	return nil
}

func isNumber(l *Lexer) bool {
	l.accept("+-")

	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}

	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}

	if r := l.current(); r != returnMeta {
		return false
	}

	l.pos--

	return true
}

func LexMap(l *Lexer) stateFn {
	return LexDetect
}

func LexConstant(l *Lexer) stateFn {
	//is it a number?
	if r := l.current(); r == '+' || r == '-' || ('0' <= r && r <= '9') {
		l.pos--
		tempPos := l.pos

		if isNumber(l) {
			l.emit(tokenNumber)
			return LexDetect
		}

		l.pos = tempPos
		l.width = 1 //because we know from the if condition
	}

	l.acceptRunUntil("\n")
	//is it a null value
	currentSlice := l.currentSlice()
	if currentSlice == nullMeta {
		l.emit(tokenNull)
	} else if currentSlice == trueMeta {
		//is it a boolean value
		l.emit(tokenTrue)
	} else if currentSlice == falseMeta {
		//is it a boolean value
		l.emit(tokenFalse)
	} else {
		if l.isSliceContain(':') {
			return LexMap
		}

		//is it a string
		l.emit(tokenString)
	}

	return LexDetect
}

func LexSpace(l *Lexer) stateFn {
	l.acceptRun(" ")
	l.emit(tokenSpace)

	return LexDetect
}

func isContainMap(l *Lexer) bool {

	return false
}

func LexDetect(l *Lexer) stateFn {
	for {
		//@
		if l.newLine && strings.HasPrefix(l.input[l.pos:], configMeta) {
			return LexConfig
		}

		//array
		if strings.HasPrefix(l.input[l.pos:], arrayMeta) {
			return LexArrayStart
		}

		l.newLine = false

		switch r := l.next(); {
		case r == eof:
			return LexEnd
		case r == returnMeta:
			return LexEnter
		case r == spaceMeta:
			return LexSpace
		default:
			return LexConstant
		}
	}
}
