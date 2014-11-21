// Copyright 2014, Ali Najafizadeh. All rights reserved.
// Use of this source code is governed by a BSD-style
// Author: Ali Najafizadeh

package note

import (
	"fmt"
	"strings"

	"github.com/alinz/gonote/note/util"
)

//Parser the parser implemenation for note
type Parser struct {
	//we are using stack since note file can be read from other file or network.
	lexers       util.Stack
	currentLexer *Lexer
	tree         *Node
}

//LoadFile load a file based on local or network
//path can be file on local or network using http
func (p *Parser) LoadFile(path string) {
	var err error
	var lexer *Lexer

	//check whether path contains http
	if lowerCasePath := strings.ToLower(path); strings.HasPrefix(lowerCasePath, "http://") ||
		strings.HasPrefix(lowerCasePath, "https://") {
		lexer, err = NewLexerWithURI(path, path)
	} else {
		lexer, err = NewLexerWithFilename(path, path)
	}

	if err != nil {
		panic(err)
	}

	p.lexers.Push(lexer)

	p.doParse()
}

func (p *Parser) doParse() {
	item, exists := p.lexers.Pop()

	if exists && item != nil {
		lexer := (item).(*Lexer)
		p.currentLexer = lexer
		for {
			if p.currentLexer == nil {
				break
			}

			tok := lexer.NextToken()

			p.process(tok)
		}
	}
}

func (p *Parser) processCommand(tok token) {
	if strings.HasPrefix(tok.val, "import ") {
		if segments := strings.Split(tok.val, " "); len(segments) == 2 {
			//push the current lexer to stack
			p.lexers.Push(p.currentLexer)
			//load the new requested file
			p.LoadFile(segments[1])
		} else {
			panic("import error: " + tok.val)
		}
	} else {
		panic("unknown command: " + tok.val)
	}
}

func (p *Parser) process(tok token) {
	switch {
	case tok.typ == tokenCommand:
		p.processCommand(tok)
	case tok.typ == tokenEnd:
		p.currentLexer = nil
	case tok.typ == tokenError:
		panic(tok)
	default:
		fmt.Println(tok)
	}
}

//Tree returns the root to parse tree
func (p *Parser) Tree() *Node {
	return nil
}

//NewParser creates a new Parser
func NewParser() *Parser {
	return &Parser{}
}
