package lexer

import "strings"

const (
	indentMeta = "\t"
	arrayMeta  = "- "
	spaceMeta  = ' '
	configMeta = "@"
	returnMeta = '\n'
	eof        = -1
)

func lexArrayStart(l *Lexer) stateFn {
	l.pos += len(arrayMeta)
	l.emit(tokenArray)
	return lexDetect
}

func lexConfig(l *Lexer) stateFn {
	l.pos += len(configMeta)
	l.start = l.pos

	l.acceptRunUntil("\n")

	l.emit(tokenCommand)

	return lexDetect
}

func lexEnter(l *Lexer) stateFn {
	l.emit(tokenEnter)
	l.newLine = true
	return lexDetect
}

func lexEnd(l *Lexer) stateFn {
	l.width = 1
	l.emit(tokenEnd)
	return nil
}

func lexMap(l *Lexer) stateFn {
	l.emit(tokenMap)

	return lexDetect
}

func lexMapOrConstant(l *Lexer) stateFn {
	l.acceptRunUntil("\n")

	//string might be containing a map or maps
	if mapPos := l.indexSlice(':'); mapPos != -1 {
		//we need to make the pos to mapPos
		l.pos -= (mapPos - 1)
		return lexMap
	}

	//is it a constant
	l.emit(tokenConstant)

	return lexDetect
}

func lexSpace(l *Lexer) stateFn {
	l.acceptRun(" ")
	l.emit(tokenSpace)

	return lexDetect
}

func lexDetect(l *Lexer) stateFn {
	for {
		//@
		if l.newLine && strings.HasPrefix(l.input[l.pos:], configMeta) {
			return lexConfig
		}

		//array
		if strings.HasPrefix(l.input[l.pos:], arrayMeta) {
			return lexArrayStart
		}

		l.newLine = false

		switch r := l.next(); {
		case r == eof:
			return lexEnd
		case r == returnMeta:
			return lexEnter
		case r == spaceMeta:
			return lexSpace
		default:
			return lexMapOrConstant
		}
	}
}
