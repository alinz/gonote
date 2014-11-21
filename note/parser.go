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
	lexers util.Stack
	tree   *Node
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
		for {
			lexer := (item).(Lexer)
			tok := lexer.NextToken()

			p.process(tok)

			if tok.typ == tokenEnd {
				break
			} else if tok.typ == tokenError {
				panic(tokenError)
			}
		}
	}
}

func (p *Parser) process(tok token) {
	fmt.Println(tok)
}

//Tree returns the root to parse tree
func (p *Parser) Tree() *Node {
	return nil
}

//NewParser creates a new Parser
func NewParser() *Parser {
	return &Parser{}
}
