package btree

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"
)

func NewBTreeNode(leaf bool) *Node {
	return &Node{
		leaf:     leaf,
		keys:     make([]string, 0, 2*degree-1),
		children: make([]*Node, 0, 2*degree),
	}
}

func NewBTree() *Tree {
	return &Tree{root: NewBTreeNode(true)}
}

// переписать на правильную вставку
// почему максимум не такой как ожидалось

func (t *Tree) Insert(key, value string) {
	root := t.root
	if len(root.keys) == 2*degree-1 {
		newRoot := NewBTreeNode(false)
		newRoot.children = append(newRoot.children, root)
		t.splitChild(newRoot, 0)
		t.root = newRoot
	}
	t.insertNonFull(t.root, key, value)
}

func (t *Tree) insertNonFull(node *Node, key, value string) {
	i := len(node.keys) - 1
	if node.leaf {
		node.keys = append(node.keys, "")
		node.data = append(node.data, "") // Добавляем место для данных
		for i >= 0 && key < node.keys[i] {
			node.keys[i+1] = node.keys[i]
			node.data[i+1] = node.data[i] // Перемещаем данные вместе с ключами
			i--
		}
		node.keys[i+1] = key
		node.data[i+1] = value // Вставляем данные
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
		t.insertNonFull(node.children[i], key, value)
	}
}

func (t *Tree) splitChild(parent *Node, index int) {
	child := parent.children[index]
	newChild := NewBTreeNode(child.leaf)

	parent.keys = append(parent.keys, "")
	parent.data = append(parent.data, "") // Добавляем место для данных
	copy(parent.keys[index+1:], parent.keys[index:])
	copy(parent.data[index+1:], parent.data[index:]) // Копируем данные
	parent.keys[index] = child.keys[degree-1]
	parent.data[index] = child.data[degree-1] // Копируем данные

	parent.children = append(parent.children, nil)
	copy(parent.children[index+2:], parent.children[index+1:])
	parent.children[index+1] = newChild

	newChild.keys = append(newChild.keys, child.keys[degree:]...)
	newChild.data = append(newChild.data, child.data[degree:]...) // Копируем данные
	child.keys = child.keys[:degree-1]
	child.data = child.data[:degree-1] // Обрезаем данные

	if !child.leaf {
		newChild.children = append(newChild.children, child.children[degree:]...)
		child.children = child.children[:degree]
	}
}

func (t *Tree) Search(key string) (string, bool) {
	return t.search(t.root, key)
}

func (t *Tree) search(node *Node, key string) (string, bool) {
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}
	if i < len(node.keys) && key == node.keys[i] {
		return node.data[i], true // Возвращаем данные
	}
	if node.leaf {
		return "", false
	}
	return t.search(node.children[i], key)
}

func (t *Tree) Delete(key string) {
	t.delete(t.root, key)
	if len(t.root.keys) == 0 && !t.root.leaf {
		t.root = t.root.children[0]
	}
}

func (t *Tree) delete(node *Node, key string) {
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

func (t *Tree) getPred(node *Node, index int) string {
	curr := node.children[index]
	for !curr.leaf {
		curr = curr.children[len(curr.keys)]
	}
	return curr.keys[len(curr.keys)-1]
}

func (t *Tree) getSucc(node *Node, index int) string {
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

	child.keys = append([]string{node.keys[index-1]}, child.keys...)
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

func CountDepth(root *Tree) int {
	return countDepth(root.root, 0)
}

func countDepth(node *Node, depth int) int {
	if node == nil {
		return depth
	}
	maxDepth := depth
	for _, child := range node.children {
		tempDepth := countDepth(child, depth+1)
		if maxDepth < tempDepth {
			maxDepth = tempDepth
		}
	}
	return maxDepth
}

func CountLoadFactorOfNode(root *Tree) (int, float64, int) {
	allLF := []int{}
	loadFactor(root.root, &allLF)
	slices.Sort(allLF)
	size := len(allLF)
	sum := 0
	for _, v := range allLF {
		sum += v
	}
	return allLF[size-1], float64(sum) / float64(size), allLF[0]
}

func loadFactor(node *Node, allLF *[]int) {
	if node == nil {
		return
	}
	*allLF = append(*allLF, len(node.keys))
	for _, child := range node.children {
		loadFactor(child, allLF)
	}
}

func LoadDataset(filename string, btree *Tree) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, record := range records {
		key := record[0] // там в датасете UUIDы в первой, не хочу парится
		keys = append(keys, key)
		value := record[1]
		btree.Insert(key, value)
	}

	return keys, nil
}
