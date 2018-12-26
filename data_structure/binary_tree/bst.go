package main

import ()

//定义节点类型
type Node struct {
	Value  int
	Parent *Node
	Left   *Node
	Right  *Node
}

//返回一个新的节点， 并定义节点的值
func NewNode(i int) *Node {
	return &Node{Value: i}
}

//比较两个节点的大小， 大于返回1， 小于返回-1， 等于返回0
func (n *Node) Compare(m *Node) int {
	if n.Value < m.Value {
		return -1
	} else if n.Value > m.Value {
		return 1
	} else {
		return 0
	}
}

//定义树的结构
type Tree struct {
	//head 表示树的顶点， 或者在寻找线路上的前一个点
	Head *Node //树的起点， 起点的指针值， 假如在寻找树的时候， 这个值就是在寻找的前一个点
	Size int
}

//返回一棵树， 而且定义了树的顶点位置， 定义了树的size
func NewTree(n *Node) *Tree {
	if n == nil {
		return &Tree{}
	}
	return &Tree{Head: n, Size: 1}
}

/*
	Insert, 首先若 树为空的话， 就把节点n设置为顶点
	然后假如继续加入其他结点， 那么就判断顶点的大小， 然后按照 左节点小右节点大的原则一直往下加下去， 形成一棵树，
	所以此函数可以递归处理，无论加入的点是左节点还是右节点， 函数会自动判断
*/
func (t *Tree) Insert(i int) {
	n := &Node{Value: i}
	if t.Head == nil { //如果这个树没有顶点的话， 那么定义n为这个树的顶点
		t.Head = n
		t.Size++
		return
	}

	h := t.Head
	for {
		if n.Compare(h) == -1 {
			if h.Left == nil {
				h.Left = n
				n.Parent = h
				break
			} else {
				h = h.Left
			}
		} else {
			if h.Right == nil {
				h.Right = n
				n.Parent = h
				break
			} else {
				h = h.Right
			}
		}
	}
	t.Size++
}

//寻找树的某个节点，返回一个 节点的指针
func (t *Tree) Search(i int) *Node {
	h := t.Head
	n := &Node{Value: i}

	for h != nil {
		switch h.Compare(n) {
		case -1:
			h = h.Right
		case 1:
			h = h.Left
		case 0:
			return h
		default:
			panic("Node not found")
		}
	}
	panic("Node not found")
}

//删除某个节点
func (t *Tree) Delete(i int) bool {
	var parent *Node
	h := t.Head
	n := &Node{Value: i}
	for h != nil {
		switch n.Compare(h) {
		case -1:
			parent = h
			h = h.Left
		case 1:
			parent = h
			h = h.Right
		case 0:
			/*删除节点逻辑*/
			if h.Left != nil {
				h.Value = h.Left.Value
				h.Left = h.Left.Left
				h.Right = h.Left.Right

				/**
				* 如果是左右节点都有的情况下,就用子左节点去代替父节点， 然后以子右节点为顶点，重新再建一颗新的树， 就是subTree
				**/
				right := h.Right
				if right != nil {
					subTree := &Tree{Head: h}
					IterOnTree(right, func(n *Node) {
						subTree.Insert(n.Value)
					})
				}
				t.Size--
				return true
			}

			if h.Right != nil {
				h.Value = h.Right.Value
				h.Left = h.Right.Left
				h.Right = h.Right.Right
				t.Size--
				return true
			}

			if parent.Left == n {
				parent.Left = nil
			} else {
				parent.Right = nil
			}
			t.Size--
			return true
		}
	}
	return false
}

func IterOnTree(n *Node, f func(*Node)) bool {
	if n == nil {
		return true
	}
	if !IterOnTree(n.Left, f) {
		return false
	}

	f(n)

	return IterOnTree(n.Right, f)
}
