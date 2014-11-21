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

//Node this is a node interface which all related nodes are inherits from this.
//we are using unexported to prevent outside package creating brand new node
type Node interface {
	Type() NodeType
	unexported()
}

//*****************************************************************************
// NodeMap
//*************************************

//NodeMap is a map node
type NodeMap struct {
	Value map[string]*Node
}

//Type returns the type of node which is NodeMapType
func (nm *NodeMap) Type() NodeType {
	return NodeMapType
}

func (nm *NodeMap) unexported() {
}

//NewNodeMap creates an empty NodeMap
func NewNodeMap() *NodeMap {
	return &NodeMap{
		Value: make(map[string]*Node),
	}
}

//*****************************************************************************
// NodeArray
//*************************************

//NodeArray is an array node
type NodeArray struct {
	Value []*Node
}

//Type returns the type of node which is NodeMapType
func (na *NodeArray) Type() NodeType {
	return NodeArrayType
}

func (na *NodeArray) unexported() {
}

//Append appends a new node to array
func (na *NodeArray) Append(node *Node) {
	na.Value = append(na.Value, node)
}

//NewNodeArray creates an empty NodeArray with cap of 5
func NewNodeArray() *NodeArray {
	return &NodeArray{
		Value: make([]*Node, 5),
	}
}

//*****************************************************************************
// NodeConstant
//*************************************

//NodeConstant is a constant node
type NodeConstant struct {
	Value string
}

//Type returns the type of node which is NodeMapType
func (nc *NodeConstant) Type() NodeType {
	return NodeConstantType
}

func (nc *NodeConstant) unexported() {
}

//NewNodeConstant creates a new node
func NewNodeConstant(value string) *NodeConstant {
	return &NodeConstant{
		Value: value,
	}
}
