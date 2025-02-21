package btree

import (
	"fmt"
	"strings"
)

func NewBTreeNode(leaf bool) *Node {
	return &Node{
		leaf:     leaf,
		keys:     make([]int, 0, 2*degree-1),
		children: make([]*Node, 0, 2*degree),
	}
}

func NewBTree() *Tree {
	return &Tree{root: NewBTreeNode(true)}
}

// todo
// dataset (любой с кегля) > btree
// load factor, depth, ноды удовлетворяют условиям B дерева
// на том же датасете посмотреть в профайлере боттлнеки

func (t *Tree) Insert(key int) {
	root := t.root
	if len(root.keys) == 2*degree-1 {
		newRoot := NewBTreeNode(false)
		newRoot.children = append(newRoot.children, root)
		t.splitChild(newRoot, 0)
		t.root = newRoot
	}
	t.insertNonFull(t.root, key)
}

func (t *Tree) insertNonFull(node *Node, key int) {
	i := len(node.keys) - 1
	if node.leaf {
		node.keys = append(node.keys, 0)
		for i >= 0 && key < node.keys[i] {
			node.keys[i+1] = node.keys[i]
			i--
		}
		node.keys[i+1] = key
	} else {
		for i >= 0 && key < node.keys[i] {
			i--
		}
		i++
		if len(node.children[i].keys) == 2*degree-1 {
			t.splitChild(node, i)
			if key > node.keys[i] {
				i++
			}
		}
		t.insertNonFull(node.children[i], key)
	}
}

func (t *Tree) splitChild(parent *Node, index int) {
	child := parent.children[index]
	newChild := NewBTreeNode(child.leaf)

	parent.keys = append(parent.keys, 0)
	copy(parent.keys[index+1:], parent.keys[index:])
	parent.keys[index] = child.keys[degree-1]

	parent.children = append(parent.children, nil)
	copy(parent.children[index+2:], parent.children[index+1:])
	parent.children[index+1] = newChild

	newChild.keys = append(newChild.keys, child.keys[degree:]...)
	child.keys = child.keys[:degree-1]

	if !child.leaf {
		newChild.children = append(newChild.children, child.children[degree:]...)
		child.children = child.children[:degree]
	}
}

func (t *Tree) Search(key int) bool {
	return t.search(t.root, key)
}

func (t *Tree) search(node *Node, key int) bool {
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}
	if i < len(node.keys) && key == node.keys[i] {
		return true
	}
	if node.leaf {
		return false
	}
	return t.search(node.children[i], key)
}

func (t *Tree) Delete(key int) {
	t.delete(t.root, key)
	if len(t.root.keys) == 0 && !t.root.leaf {
		t.root = t.root.children[0]
	}
}

func (t *Tree) delete(node *Node, key int) {
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}

	if i < len(node.keys) && key == node.keys[i] {
		if node.leaf {
			t.deleteFromLeaf(node, i)
		} else {
			t.deleteFromInternalNode(node, i)
		}
	} else {
		if node.leaf {
			fmt.Println("Key not found:", key)
			return
		}
		flag := (i == len(node.keys))
		if len(node.children[i].keys) < degree {
			t.fill(node, i)
		}
		if flag && i > len(node.keys) {
			t.delete(node.children[i-1], key)
		} else {
			t.delete(node.children[i], key)
		}
	}
}

func (t *Tree) deleteFromLeaf(node *Node, index int) {
	node.keys = append(node.keys[:index], node.keys[index+1:]...)
}

func (t *Tree) deleteFromInternalNode(node *Node, index int) {
	key := node.keys[index]
	if len(node.children[index].keys) >= degree {
		pred := t.getPred(node, index)
		node.keys[index] = pred
		t.delete(node.children[index], pred)
	} else if len(node.children[index+1].keys) >= degree {
		succ := t.getSucc(node, index)
		node.keys[index] = succ
		t.delete(node.children[index+1], succ)
	} else {
		t.merge(node, index)
		t.delete(node.children[index], key)
	}
}

func (t *Tree) getPred(node *Node, index int) int {
	curr := node.children[index]
	for !curr.leaf {
		curr = curr.children[len(curr.keys)]
	}
	return curr.keys[len(curr.keys)-1]
}

func (t *Tree) getSucc(node *Node, index int) int {
	curr := node.children[index+1]
	for !curr.leaf {
		curr = curr.children[0]
	}
	return curr.keys[0]
}

func (t *Tree) fill(node *Node, index int) {
	if index != 0 && len(node.children[index-1].keys) >= degree {
		t.borrowFromPrev(node, index)
	} else if index != len(node.keys) && len(node.children[index+1].keys) >= degree {
		t.borrowFromNext(node, index)
	} else {
		if index != len(node.keys) {
			t.merge(node, index)
		} else {
			t.merge(node, index-1)
		}
	}
}

func (t *Tree) borrowFromPrev(node *Node, index int) {
	child := node.children[index]
	sibling := node.children[index-1]

	child.keys = append([]int{node.keys[index-1]}, child.keys...)
	node.keys[index-1] = sibling.keys[len(sibling.keys)-1]

	if !child.leaf {
		child.children = append([]*Node{sibling.children[len(sibling.children)-1]}, child.children...)
		sibling.children = sibling.children[:len(sibling.children)-1]
	}

	sibling.keys = sibling.keys[:len(sibling.keys)-1]
}

func (t *Tree) borrowFromNext(node *Node, index int) {
	child := node.children[index]
	sibling := node.children[index+1]

	child.keys = append(child.keys, node.keys[index])
	node.keys[index] = sibling.keys[0]
	sibling.keys = sibling.keys[1:]

	if !child.leaf {
		child.children = append(child.children, sibling.children[0])
		sibling.children = sibling.children[1:]
	}
}

func (t *Tree) merge(node *Node, index int) {
	child := node.children[index]
	sibling := node.children[index+1]

	child.keys = append(child.keys, node.keys[index])
	child.keys = append(child.keys, sibling.keys...)

	if !child.leaf {
		child.children = append(child.children, sibling.children...)
	}

	node.keys = append(node.keys[:index], node.keys[index+1:]...)
	node.children = append(node.children[:index+1], node.children[index+2:]...)
}

func (t *Tree) Print() {
	t.printTree(t.root, 0)
}

func (t *Tree) printTree(node *Node, level int) {
	if node == nil {
		return
	}
	fmt.Printf("Level %d: %v\n", level, node.keys)
	for _, child := range node.children {
		t.printTree(child, level+1)
	}
}

func (t *Tree) PrettyPrint() {
	t.prettyPrintTree(t.root, 0)
}

func (t *Tree) prettyPrintTree(node *Node, level int) {
	if node == nil {
		return
	}
	var sb strings.Builder
	sb.WriteString(strings.Repeat("  ", level))
	sb.WriteString(fmt.Sprintf("Level %d: %v\n", level, node.keys))
	fmt.Print(sb.String())
	for _, child := range node.children {
		t.prettyPrintTree(child, level+1)
	}
}
