package main

import (
	"errors"
	"fmt"
)

type List struct { //定义链表结构
	Length int   //链表长度
	Head   *Node //链表头
	Tail   *Node //链表尾
}

//新建一个链表
func NewList() *List {
	l := new(List) //新建一个结构体存在内存里面， 用new关键字
	l.Length = 0
	return l
}

type Node struct { //node 结构体
	Value interface{} //Node的值
	Prev  *Node       //前一个node
	Next  *Node       //后一个node
}

//新建一个node
func NewNode(value interface{}) *Node {
	return &Node{Value: value}
}

//返回链表长度
func (l *List) Len() int {
	return l.Length
}

//判断链表是否为空
func (l *List) IsEmpty() bool {
	return l.Length == 0
}

//链表的 头点 插入法， 即从第一个节点的位置插入链表
func (l *List) Prepend(value interface{}) {
	node := NewNode(value)
	//当链表为空时， 初始化head， 然后首尾相连
	if l.Len() == 0 {
		l.Head = node
		l.Tail = l.Head
	} else {
		//当有一个新的node要进链表的时候， 原head node 就成为新的node 的next node，
		//新的node 就代替他成为head
		formerHead := l.Head
		formerHead.Prev = node

		node.Next = formerHead
		l.Head = node
	}
	l.Length++
}

//链表的 尾点 插入法， 即从最后一个节点的位置插入链表
func (l *List) Append(value interface{}) {
	node := NewNode(value)
	if l.Len() == 0 {
		l.Head = node
		l.Tail = l.Head
	} else {
		formerTail := l.Tail
		formerTail.Next = node
		node.Prev = formerTail
		l.Tail = node
	}
	l.Length++
}

//在指定的 index 位置加入某个节点
func (l *List) Add(value interface{}, index int) error {
	if index > l.Len() {
		return errors.New("Index out of range")
	}
	node := NewNode(value)
	//当链表为空， 即add位置在开头， 则使用Prepend效率高
	if l.Len() == 0 || index == 0 {
		l.Prepend(value)
		return nil
	}

	//当链表不为空， 但要add的位置在链表末尾， 则使用Append效率高
	if l.Len()-1 == index {
		l.Append(value)
		return nil
	}

	nextNode, _ := l.Get(index + 1)
	prevNode := nextNode.Prev

	prevNode.Next = node
	node.Prev = prevNode

	nextNode.Prev = node
	node.Next = nextNode

	l.Length++
	return nil
}

/*
** 返回error类型
** 移除链表里某个元素
 */
func (l *List) Remove(value interface{}) error {
	if l.Len() == 0 {
		return errors.New("Empty list")
	}
	//若刚好要删除的是第一个节点
	if l.Head.Value == value {
		l.Head = l.Head.Next
		l.Length--
		return nil
	}
	found := 0
	//循环链表
	for n := l.Head; n != nil; n = n.Next {
		if *n.Value.(*Node) == value && found == 0 {
			n.Next.Prev, n.Prev.Next = n.Prev, n.Next
			l.Length--
			found++
		}
	}
	if found == 0 {
		return errors.New("Node not found")
	}
	return nil
}

//获取某个节点
//返回 node 和 error
//都是从第一个节点开始获取， 时间O(n)
func (l *List) Get(index int) (*Node, error) {
	if index > l.Len() {
		return nil, errors.New("index out of range")
	}
	node := l.Head //从第一个节点开始获取
	for i := 0; i < index; i++ {
		node = node.Next
	}
	return node, nil
}

//合并两个链表
func (l *List) Concat(k *List) {
	l.Tail.Next, k.Tail.Prev = k.Head, l.Tail
	l.Tail = k.Tail
	l.Length += k.Length
}

func (list *List) Map(f func(node *Node)) {
	for node := list.Head; node != nil; node = node.Next {
		n := node.Value.(*Node)
		f(n)
	}
}

func (list *List) Each(f func(node Node)) {
	for node := list.Head; node != nil; node = node.Next {
		f(*node)
	}
}

//在链表中寻找某个节点， 并返回index位置
func (l *List) Find(node *Node) (int, error) {
	if l.Len() == 0 {
		return 0, errors.New("Empty list")
	}

	index := 0
	found := -1
	l.Map(func(n *Node) {
		index++
		if n.Value == node.Value && found == -1 {
			found = index
		}
	})
	if found == -1 {
		return 0, errors.New("Item not found")
	}
	return found, nil
}

func main() {
	l := NewList()
	l.Prepend(NewNode(1))
	l.Prepend(NewNode(2))
	l.Prepend(NewNode(3))

	zero := *slice(l.Get(0))[0].(*Node).Value.(*Node)
	one := *slice(l.Get(1))[0].(*Node).Value.(*Node)
	two := *slice(l.Get(2))[0].(*Node).Value.(*Node)

	if zero != *NewNode(3) || one != *NewNode(2) || two != *NewNode(1) {
		fmt.Println(*one.Value.(*Node), *NewNode(2))
		fmt.Println(zero.Value)
	}

	fmt.Println(zero.Value)
}

func slice(args ...interface{}) []interface{} {
	return args
}
