package main

import (
	. "./base"
	"fmt"
)

type CNode struct {
	data Object
	next *CNode
}

type CList struct {
	size uint64
	head *CNode
}

func (cList *CList) Init() {
	lst := *cList
	lst.size = 0
	lst.head = nil
}

func (cList *CList) GetSize() uint64 {
	return (*cList).size
}

func (cList *CList) GetHead() *CNode {
	return (*cList).head
}

func (node *CNode) GetData() Object {
	return (*node).data
}

func (node *CNode) GetNext() *CNode {
	return (*node).next
}

func (cList *CList) Append(data Object) bool {
	node := new(CNode)
	(*node).data = data
	if cList.GetSize() == 0 {
		(*cList).head = node
	} else {
		item := cList.GetHead()
		for ; (*item).next != cList.GetHead(); item = (*item).next {
		}
		(*item).next = node
	}

	(*node).next = (*cList).head
	(*cList).size++

	return true
}

func (cList *CList) InsertNext(elmt *CNode, data Object) bool {
	if elmt == nil {
		return false
	}
	node := new(CNode)
	(*node).data = data

	(*node).next = (*elmt).next
	(*elmt).next = node

	(cList).size++
	return true
}

func (cList *CList) Romove(elmt *CNode) Object {
	if elmt == nil {
		return false
	}
	item := cList.GetHead()
	for ; (*item).next != elmt; item = (*item).next {
	}
	(*item).next = (*elmt).next
	(*cList).size--
	return elmt.GetData()
}

func main() {
	l := new(CList)
	l.Append(1)
	l.Append(2)
	l.Append(3)

	/*fmt.Println(l.GetSize())
	fmt.Println(l.GetHead().GetData())*/

	l.InsertNext(l.GetHead(), 7)

	fmt.Println(l.GetSize())
	fmt.Println(l.GetHead())
	fmt.Println(l.head.GetNext())

	l.Romove(l.head.next)
	fmt.Println(l.head.next)

	fmt.Println(l.size)
}
