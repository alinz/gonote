// Copyright 2014, Ali Najafizadeh. All rights reserved.
// Use of this source code is governed by a BSD-style
// Author: Ali Najafizadeh

package note

import (
	"strconv"
	"strings"
)

const (
	arrayMeta   = "-"
	spaceMeta   = ' '
	mapMeta     = ':'
	commandMeta = "@"
	returnMeta  = '\n'
	eof         = -1
)

func lexArrayStart(l *Lexer) stateFn {
	l.pos += len(arrayMeta)
	l.emit(tokenArray)

	//we are ignoring spaces after an array is found
	l.acceptRun(strconv.QuoteRuneToASCII(spaceMeta))
	l.ignore()

	l.newArray = true

	return lexDetect
}

func lexCommand(l *Lexer) stateFn {
	l.pos += len(commandMeta)
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
	//these two increament skip the `:` chars
	l.pos++
	l.start++

	//ignore spaces
	//for example
	//name:     john
	//we don't want to parse spaces between `:` and `john`
	l.acceptRun(strconv.QuoteRuneToASCII(spaceMeta))
	l.ignore()

	l.newMap = true

	return lexDetect
}

func lexMapOrConstant(l *Lexer) stateFn {
	l.acceptRunUntil("\n")

	//string might be containing a map or maps
	if mapPos := l.indexSlice(mapMeta); mapPos != -1 {
		//we need since mapPos is the index from slice,
		//we need to add to add it to l.start to find the new pos
		l.pos = l.start + mapPos
		return lexMap
	}

	//is it a constant
	l.emit(tokenConstant)

	return lexDetect
}

func lexSpace(l *Lexer) stateFn {
	l.acceptRun(strconv.QuoteRuneToASCII(spaceMeta))
	l.emit(tokenSpace)

	return lexDetect
}

func lexDetect(l *Lexer) stateFn {
	for {
		//command
		if (l.newLine || l.newArray || l.newMap) &&
			strings.HasPrefix(l.input[l.pos:], commandMeta) {
			return lexCommand
		}

		//array
		if strings.HasPrefix(l.input[l.pos:], arrayMeta) {
			return lexArrayStart
		}

		l.newLine = false
		l.newArray = false
		l.newMap = false

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
