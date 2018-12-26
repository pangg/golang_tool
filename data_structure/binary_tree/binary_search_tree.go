/*
	type BiSearchTree struct
	func (bst *BiSearchTree) Add(data float64)                        //插入节点
	func (bst *BiSearchTree) Delete(data float64) 	                  //删除节点
	func (bst BiSearchTree) GetRoot() *TreeNode                       //获取根节点
	func (bst BiSearchTree) IsEmpty() bool                            //检查树是否为空
	func (bst BiSearchTree) InOrderTravel()                           //中序遍历(也就是从小到大输出)
	func (bst BiSearchTree) Search(data float64) *TreeNode            //查找节点
	func (bst BiSearchTree) GetDeepth() int                           //获取树的深度
	func (bst BiSearchTree) GetMin() float64                          //获取值最小的节点
	func (bst BiSearchTree) GetMax() float64                          //获取值最大的节点
	func (bst BiSearchTree) GetPredecessor(data float64) *TreeNode    //获取直接前驱
	func (bst BiSearchTree) GetSuccessor(data float64) *TreeNode      //获取直接后继
	func (bst *BiSearchTree) Clear()                                  //清空树

*/

package main

import (
	"fmt"
)

type TreeNode struct {
	data   float64
	lchild *TreeNode
	rchild *TreeNode
	parent *TreeNode
}

type BiSearchTree struct {
	root   *TreeNode
	cur    *TreeNode
	create *TreeNode
}

func (bst *BiSearchTree) Add(data float64) {
	bst.create = new(TreeNode)
	bst.create.data = data

	if !bst.IsEmpty() {
		bst.cur = bst.root
		for {
			if data < bst.cur.data {
				//如果要插入的值比当前节点的值小，则当前节点指向当前节点的左孩子，如果
				//左孩子为空，就在这个左孩子上插入新值
				if bst.cur.lchild == nil {
					bst.cur.lchild = bst.create
					bst.create.parent = bst.cur
					break
				} else {
					bst.cur = bst.cur.lchild
				}
			} else if data > bst.cur.data {
				//如果要插入的值比当前节点的值大，则当前节点指向当前节点的右孩子，如果
				//右孩子为空，就在这个右孩子上插入新值
				if bst.cur.rchild == nil {
					bst.cur.rchild = bst.create
					bst.create.parent = bst.cur
					break
				} else {
					bst.cur = bst.cur.rchild
				}
			} else {
				return
			}
		}
	} else {
		bst.root = bst.create
		bst.root.parent = nil
	}
}

func (bst *BiSearchTree) Delete(data float64) {
	var (
		deleteNode func(node *TreeNode)
		node       *TreeNode = bst.Search(data)
	)
}

func (bst *BiSearchTree) GetRoot() *TreeNode {
	if bst.root != nil {
		return bst.root
	}
	return nil
}

func (bst *BiSearchTree) IsEmpty() bool {
	if bst.root == nil {
		return true
	}
	return false
}

func (bst *BiSearchTree) InOrderTravel() {
	var inOrderTravel func(node *TreeNode)
	inOrderTravel = func(node *TreeNode) {
		if node != nil {
			inOrderTravel(node.lchild)
			fmt.Printf("%g", node.data)
			inOrderTravel(node.rchild)
		}
	}

	inOrderTravel(bst.root)
}

func (bst *BiSearchTree) Search(data float64) *TreeNode {
	//和Add操作类似，只要按照比当前节点小就往左孩子上拐，比当前节点大就往右孩子上拐的思路
	//一路找下去，知道找到要查找的值返回即可
	bst.cur = bst.root
	for {
		if bst.cur == nil {
			return nil
		}
		if data < bst.cur.data {
			bst.cur = bst.cur.lchild
		} else if data > bst.cur.data {
			bst.cur = bst.cur.rchild
		} else {
			return bst.cur
		}
	}
}

func (bst *BiSearchTree) GetDeepth() int {
	var getDeepth func(node *TreeNode) int
	getDeepth = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		if node.lchild == nil && node.rchild == nil {
			return 1
		}
		var (
			ldeepth int = getDeepth(node.lchild)
			rdeepth int = getDeepth(node.rchild)
		)
		if ldeepth > rdeepth {
			return ldeepth + 1
		} else {
			return rdeepth + 1
		}
	}

	return getDeepth(bst.root)
}

func (bst *BiSearchTree) GetMin() float64 {
	//根据二叉查找树的性质，树中最左边的节点就是值最小的节点
	if bst.root == nil {
		return -1
	}
	bst.cur = bst.root
	for {
		if bst.cur.lchild != nil {
			bst.cur = bst.cur.lchild
		} else {
			return bst.cur.data
		}
	}
}

func (bst *BiSearchTree) GetMax() float64 {
	//根据二叉查找树的性质，树中最右边的节点就是值最大的节点
	if bst.root == nil {
		return -1
	}
	bst.cur = bst.root
	for {
		if bst.cur.rchild != nil {
			bst.cur = bst.cur.rchild
		} else {
			return bst.cur.data
		}
	}
}
