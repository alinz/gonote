// Copyright 2014, Ali Najafizadeh. All rights reserved.
// Use of this source code is governed by a BSD-style
// Author: Ali Najafizadeh

package note

import "github.com/alinz/gonote/note/util"

//Parser the parser implemenation for note
type Parser struct {
	//we are using stack since note file can be read from other file or network.
	lexers util.Stack
}

//LoadFile load a file based on local or network
func (p *Parser) LoadFile(path string) {

}

//NewParser creates a new Parser
func NewParser() *Parser {
	return &Parser{}
}
