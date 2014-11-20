// Copyright 2014, Ali Najafizadeh. All rights reserved.
// Use of this source code is governed by a BSD-style
// Author: Ali Najafizadeh

package note

import "fmt"

type tokenType int

const (
	//DocStart when a new document is loaded or parsed
	tokenDocStart tokenType = iota
	tokenArray
	tokenMap
	tokenConstant
	tokenSpace
	tokenEnter
	tokenCommand
	tokenError
	tokenEnd
)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	tokType := ""
	switch t.typ {
	case tokenEnter:
		tokType = "tokenEnter"
	case tokenCommand:
		tokType = "tokenCommand"
	case tokenError:
		tokType = "tokenError"
	case tokenArray:
		tokType = "tokenArray"
	case tokenMap:
		tokType = "tokenMap"
	case tokenConstant:
		tokType = "tokenConstant"
	case tokenSpace:
		tokType = "tokenSpace"
	case tokenEnd:
		tokType = "tokenEnd"
	}

	if len(t.val) > 20 {
		return fmt.Sprintf("%s: %.20q...", tokType, t.val)
	}
	return fmt.Sprintf("%s: %q", tokType, t.val)
}
