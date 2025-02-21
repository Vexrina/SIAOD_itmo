package main

import (
	"fmt"
	"vexrina/siaod_itmo/lab_02/btree"
)

func main() {
	head := btree.NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, key := range keys {
		head.Insert(key)
	}

	fmt.Println("B-Tree after insertion:")
	head.PrettyPrint()

	fmt.Println("Search for 6:", head.Search(6))   // true
	fmt.Println("Search for 15:", head.Search(15)) // false

	head.Delete(6)
	fmt.Println("B-Tree after deleting 6:")
	head.PrettyPrint()

	head.Delete(12)
	fmt.Println("B-Tree after deleting 12:")
	head.PrettyPrint()
}
