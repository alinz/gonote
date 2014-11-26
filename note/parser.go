// Copyright 2014, Ali Najafizadeh. All rights reserved.
// Use of this source code is governed by a BSD-style
// Author: Ali Najafizadeh

package note

import (
	"fmt"
	"strings"

	"github.com/alinz/gonote/note/util"
)

type nodeExtra struct {
	node        Node
	indentation int
}

//Parser the parser implemenation for note
type Parser struct {
	//we are using stack since note file can be read from other file or network.
	lexers       util.Stack
	currentLexer *Lexer     //holds the current active lexer
	tree         Node       //holds the root of parse tree
	current      Node       //holds the pointer of current node
	indentation  int        //keeps track of current indentation
	nodes        util.Stack //stores processed node in stack
	lastKeyMap   string
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
	log(tok)
	switch tok.typ {
	case tokenArray:
		p.indentation += len(tok.val)
		p.currentNode(NodeArrayType)

	case tokenMap:
		p.lastKeyMap = tok.val
		p.currentNode(NodeMapType)

	case tokenConstant:
		node := NewNodeConstant(tok.val)
		p.addNodeToCurrent(node)

	case tokenEnter:
		p.indentation = 0

	case tokenSpace:
		p.indentation += len(tok.val)

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

func (p *Parser) currentNode(nodeType NodeType) {
	if p.current == nil {
		node := p.makeNode(nodeType)
		p.current = node
		p.nodes.Push(&nodeExtra{
			node:        node,
			indentation: p.indentation,
		})
		p.setupRoot(node)
	} else {
		for {
			temp, exists := p.nodes.Pop()
			if !exists {
				p.current = nil
				break
			}

			if ptr := (temp).(*nodeExtra); ptr.node.Type() == nodeType && ptr.indentation == p.indentation {
				p.current = ptr.node
				p.nodes.Push(temp)
				break
			}
		}
	}
}

func (p *Parser) makeNode(nodeType NodeType) (node Node) {
	switch nodeType {
	case NodeArrayType:
		node = NewNodeArray()
	case NodeMapType:
		node = NewNodeMap()
	case NodeConstantType:
		panic("can not create a complex node")
	}
	return
}

func (p *Parser) addNodeToCurrent(node Node) {
	if p.tree == nil {
		p.tree = node
		p.current = node
	} else if p.current != nil {

		switch p.current.Type() {
		case NodeArrayType:
			((p.current).(*NodeArray)).Append(node)
		case NodeMapType:
			((p.current).(*NodeMap)).Put(p.lastKeyMap, node)
		case NodeConstantType:
			panic("can not change/add the constant node.")
		}

	} else {
		panic("current pointer is null")
	}
}

func (p *Parser) setupRoot(node Node) {
	if p.tree == nil {
		p.tree = node
	}
}

//Tree returns the root to parse tree
func (p *Parser) Tree() Node {
	return p.tree
}

//NewParser creates a new Parser
func NewParser() *Parser {
	return &Parser{}
}
