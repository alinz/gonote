// Copyright 2014, Ali Najafizadeh. All rights reserved.
// Use of this source code is governed by a BSD-style
// Author: Ali Najafizadeh

package note

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alinz/gonote/note/util"
)

//Parser the parser implemenation for note
type Parser struct {
	//we are using stack since note file can be read from other file or network.
	lexers       util.Stack
	currentLexer *Lexer       //holds the current active lexer
	tree         Node         //holds the root of parse tree
	currentNode  Node         //holds the pointer of current node
	indentation  int          //keeps track of current indentation
	nodeIndexMap map[int]Node //the key holds the number of indentation
}

//LoadFile load a file based on local or network
func (p *Parser) LoadFile(path string) {
	var err error
	var lexer *Lexer

	lexer, err = NewLexerWithPath(path, path)

	if err != nil {
		panic(err)
	}

	p.lexers.Push(lexer)

	p.doParse()
}

func (p *Parser) doParse() {
	for {
		item, exists := p.lexers.Pop()

		if !exists || item == nil {
			break
		}

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
	case tok.typ == tokenArray:
		if node, err := p.getCurrentNode(NodeArrayType); err == nil {
			p.currentNode = node
		} else {
			panic(err)
		}

	case tok.typ == tokenConstant:
		a := (p.currentNode).(*NodeArray)
		a.Append(NewNodeConstant(tok.val))
		//log((NewNodeConstant(tok.val))

	case tok.typ == tokenEnter:
		p.indentation = 0
	case tok.typ == tokenSpace:
		p.indentation += len(tok.val)

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

func (p *Parser) getCurrentNode(nodeType NodeType) (node Node, err error) {
	err = nil
	node, ok := p.nodeIndexMap[p.indentation]
	if ok {
		if node.Type() != nodeType {
			node = nil
			err = errors.New("wrong indentation object")
		}
	} else {

		switch nodeType {
		case NodeMapType:
			node = NewNodeMap()
		case NodeArrayType:
			log("create a new array")
			node = NewNodeArray()
		default:
			err = errors.New("node can not be created as a constant node")
		}

		if err == nil {
			p.nodeIndexMap[p.indentation] = node
		}
	}

	if p.tree == nil {
		p.tree = node
	}

	return
}

//Tree returns the root to parse tree
func (p *Parser) Tree() Node {
	return p.tree
}

//NewParser creates a new Parser
func NewParser() *Parser {
	return &Parser{
		nodeIndexMap: make(map[int]Node),
	}
}
