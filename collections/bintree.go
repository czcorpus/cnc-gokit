// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collections

import "fmt"

type Comparable interface {
	// Compare should return:
	// *  x > 0 if this item is greater than the `other`,
	// *  x == 0 if items are equal
	// *  x < 0 if this item is lesser than the `other`
	Compare(other Comparable) int
}

type node[T Comparable] struct {
	lft    *node[T]
	rgt    *node[T]
	parent *node[T]
	value  Comparable
}

func (node *node[T]) isLeaf() bool {
	return node.lft == nil && node.rgt == nil
}

// BinTree is a simple unbalanced binary tree
// implementation for storing sorted values
type BinTree[T Comparable] struct {
	root   *node[T]
	length int
}

func (bt *BinTree[T]) Add(v ...T) {
	for _, vx := range v {
		bt.add(vx)
	}
}

func (bt *BinTree[T]) add(v T) {
	defer func() { bt.length++ }()
	if bt.root == nil {
		bt.root = &node[T]{value: v}
		return
	}
	currNode := bt.root
	for currNode != nil {
		cmp := v.Compare(currNode.value)
		if cmp <= 0 {
			if currNode.lft == nil {
				currNode.lft = &node[T]{value: v, parent: currNode}
				return
			}
			currNode = currNode.lft

		} else {
			if currNode.rgt == nil {
				currNode.rgt = &node[T]{value: v, parent: currNode}
				return
			}
			currNode = currNode.rgt
		}
	}
}

func (BinTree[T]) goLeftmost(root *node[T], stack []*node[T]) []*node[T] {
	curr := root.lft
	for curr != nil {
		stack = append(stack, curr)
		curr = curr.lft
	}
	return stack
}

// findNodeAt returns node at position `idx` along with its parent
func (bt *BinTree[T]) findNodeAt(idx int) *node[T] {
	stack := []*node[T]{bt.root}
	stack = bt.goLeftmost(bt.root, stack)
	var i int
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if i == idx {
			return node
		}
		i++
		if node.rgt != nil {
			stack = append(stack, node.rgt)
			stack = bt.goLeftmost(node.rgt, stack)
		}
	}
	return nil
}

func (bt BinTree[T]) ToSlice() []T {
	if bt.length == 0 {
		return []T{}
	}
	ans := make([]T, 0, bt.length)
	stack := []*node[T]{bt.root}
	stack = bt.goLeftmost(bt.root, stack)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		ans = append(ans, node.value.(T))
		if node.rgt != nil {
			stack = append(stack, node.rgt)
			stack = bt.goLeftmost(node.rgt, stack)
		}
	}
	return ans
}

func (bt *BinTree[T]) Remove(idx int) T {
	if bt.length == 0 {
		panic(fmt.Sprintf("BinTree index overflow: %d (len: %d)", idx, bt.length))
	}
	srch := bt.findNodeAt(idx)
	if srch == nil {
		var zeroVal T
		return zeroVal
	}
	rtrn := srch.value.(T)
	defer func() { bt.length-- }()
	if srch.isLeaf() {
		if srch.parent != nil {
			if srch.parent.lft == srch {
				srch.parent.lft = nil

			} else {
				srch.parent.rgt = nil
			}
		}

	} else if srch.lft == nil && srch.rgt != nil {
		if srch.parent != nil {
			if srch.parent.lft == srch {
				srch.parent.lft = srch.rgt
				srch.rgt.parent = srch.parent

			} else {
				srch.parent.rgt = srch.rgt
				srch.rgt.parent = srch.parent
			}
		}

	} else if srch.lft != nil && srch.rgt == nil {
		if srch.parent != nil {
			if srch.parent.lft == srch {
				srch.parent.lft = srch.lft
				srch.lft.parent = srch.parent

			} else {
				srch.parent.rgt = srch.lft
				srch.lft.parent = srch.parent
			}
		}

	} else {
		succ := srch.rgt
		succParent := srch
		for succ.lft != nil {
			succParent = succ
			succ = succ.lft
		}
		if succParent != srch {
			succParent.lft = succ.rgt

		} else {
			succParent.rgt = succ.rgt
		}
		srch.value = succ.value
	}
	return rtrn
}

// Get returns an item on i-th index. The function
// also supports negative indices where -1 is the last
// item.
// In case the index does not exist in data the function
// panics.
func (bt *BinTree[T]) Get(idx int) T {
	if bt.length == 0 {
		panic(fmt.Sprintf("BinTree index overflow: %d (len: %d)", idx, bt.length))
	}
	rIdx := idx
	if idx < 0 {
		rIdx = bt.length + idx
	}
	srch := bt.findNodeAt(rIdx)
	if srch != nil {
		ans, _ := srch.value.(T)
		return ans
	}
	panic(fmt.Sprintf("BinTree index overflow: %d (len: %d)", idx, bt.length))
}

func (bt *BinTree[T]) Len() int {
	return bt.length
}
