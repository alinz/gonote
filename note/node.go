// Copyright 2014, Ali Najafizadeh. All rights reserved.
// Use of this source code is governed by a BSD-style
// Author: Ali Najafizadeh

package note

//NodeType is a type that defines all node types
type NodeType int

const (
	//NodeMapType is the type of map node
	NodeMapType NodeType = iota
	//NodeArrayType is the type of array node
	NodeArrayType
	//NodeConstantType is the type of constant node
	NodeConstantType
)

//Node this is a node interface
//we are using unexported to prevent outside package creating brand new node
type Node interface {
	Type() NodeType
	unexported()
}

//NodeMap is a map node
type NodeMap struct{}

//Type returns the type of node which is NodeMapType
func (nm *NodeMap) Type() NodeType {
	return NodeMapType
}

func (nm *NodeMap) unexported() {
}

//NodeArray is an array node
type NodeArray struct{}

//Type returns the type of node which is NodeMapType
func (na *NodeArray) Type() NodeType {
	return NodeArrayType
}

func (na *NodeArray) unexported() {
}

//NodeConstant is a constant node
type NodeConstant struct{}

//Type returns the type of node which is NodeMapType
func (nc *NodeConstant) Type() NodeType {
	return NodeConstantType
}

func (nc *NodeConstant) unexported() {
}
