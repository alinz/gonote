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
	current      Node         //holds the pointer of current node
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
	switch tok.typ {
	case tokenArray:
		p.indentation += len(tok.val)
		if err := p.currentNode(NodeArrayType); err != nil {
			panic(err)
		}
	case tokenMap:
		//we have to delete the map object from indentation once it's done. somehow!
		if err := p.currentNode(NodeMapType); err != nil {
			panic(err)
		}

	case tokenConstant:
		a := (p.current).(*NodeArray)
		a.Append(NewNodeConstant(tok.val))

	case tokenEnter:
		p.indentation = 0
	case tokenSpace:
		p.indentation += len(tok.val)

	//Need to rethink about the below operation
	//
	//
	case tokenCommand:
		p.processCommand(tok)
	case tokenEnd:
		p.currentLexer = nil
	case tokenError:
		panic(tok)
	default:
		fmt.Println(tok)
	}
}

func (p *Parser) currentNode(nodeType NodeType) (err error) {
	err = nil

	//get the node from indentation
	//it can be nil
	node, ok := p.nodeIndexMap[p.indentation]

	log(node)

	if ok {
		if node.Type() != nodeType {
			err = errors.New("wrong indentation object")
		}
	} else {
		switch nodeType {
		case NodeArrayType:
			node = NewNodeArray()
			p.addToNode(node)
		case NodeMapType:
			node = NewNodeMap()
		default:
			err = errors.New("current node can not be a constant node")
		}
	}

	if p.tree == nil {
		p.tree = node
	}

	p.nodeIndexMap[p.indentation] = node

	p.current = node

	return
}

func (p *Parser) addToNode(node Node) {
	if p.current != nil {
		switch p.current.Type() {
		case NodeArrayType:
			(p.current).(*NodeArray).Append(node)
		case NodeMapType:
			panic("not implemented yet")
		default:
			panic("something went seriously wrong")
		}
	}
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
