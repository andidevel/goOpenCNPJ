//    Copyright 2021 Anderson Rodrigues do Livramento

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package importer

import "sync"

// Binary Search Tree methods

type BTreeNode struct {
	ID    string
	Data  map[string]interface{}
	Right *BTreeNode
	Left  *BTreeNode
}

type NodeVisitFunc func(*BTreeNode)

type BinarySearchTree struct {
	root *BTreeNode
	lock sync.RWMutex
}

// #### BTreeNode ####
func (n *BTreeNode) add(node *BTreeNode) {
	if n == nil {
		return
	}
	if node.ID < n.ID {
		if n.Left == nil {
			n.Left = node
			return
		}
		n.Left.add(node)
		return
	}
	if node.ID > n.ID {
		if n.Right == nil {
			n.Right = node
			return
		}
		n.Right.add(node)
	}
}

func (n *BTreeNode) find(id string) *BTreeNode {
	if n == nil {
		return nil
	}
	if id < n.ID {
		return n.Left.find(id)
	}
	if id > n.ID {
		return n.Right.find(id)
	}
	return n
}

func (n *BTreeNode) max() *BTreeNode {
	if n == nil {
		return nil
	}
	for {
		if n.Right == nil {
			return n
		}
		n = n.Right
	}
}

func (n *BTreeNode) min() *BTreeNode {
	if n == nil {
		return nil
	}
	for {
		if n.Left == nil {
			return n
		}
		n = n.Left
	}
}

func (n *BTreeNode) remove(id string) *BTreeNode {
	if n == nil {
		return nil
	}
	if id < n.ID {
		n.Left = n.Left.remove(id)
		return n
	}
	if id > n.ID {
		n.Right = n.Right.remove(id)
		return n
	}
	if n.Left == nil && n.Right == nil {
		n = nil
		return nil
	}
	if n.Left == nil {
		n = n.Right
		return n
	}
	if n.Right == nil {
		n = n.Left
		return n
	}
	minRight := n.Right.min()
	n.ID, n.Data = minRight.ID, minRight.Data
	n.Right = n.Right.remove(n.ID)
	return n
}

func (n *BTreeNode) walkInOrderTraversal(visit NodeVisitFunc) {
	if n != nil {
		n.Left.walkInOrderTraversal(visit)
		visit(n)
		n.Right.walkInOrderTraversal(visit)
	}
}

// #### BinarySearchTree ####
func (bst *BinarySearchTree) Add(id string, data map[string]interface{}) {
	bst.lock.Lock()
	defer bst.lock.Unlock()

	node := &BTreeNode{
		ID:    id,
		Data:  data,
		Right: nil,
		Left:  nil,
	}
	if bst.root == nil {
		bst.root = node
		return
	}
	bst.root.add(node)
}

func (bst *BinarySearchTree) Min() *BTreeNode {
	bst.lock.RLock()
	defer bst.lock.RUnlock()

	return bst.root.min()
}

func (bst *BinarySearchTree) Max() *BTreeNode {
	bst.lock.RLock()
	defer bst.lock.RUnlock()

	return bst.root.max()
}

func (bst *BinarySearchTree) Find(id string) *BTreeNode {
	bst.lock.RLock()
	defer bst.lock.RUnlock()

	return bst.root.find(id)
}

func (bst *BinarySearchTree) Remove(id string) {
	bst.lock.Lock()
	defer bst.lock.Unlock()

	if bst.root == nil {
		return
	}
	bst.root.remove(id)
}

func (bst *BinarySearchTree) Walk(visit NodeVisitFunc) {
	bst.lock.Lock()
	defer bst.lock.Unlock()

	bst.root.walkInOrderTraversal(visit)
}
