package main

import (
	"fmt"
	"vexrina/siaod_itmo/lab_02/btree"
)

func main() {
	head := btree.NewBTree()
	head.Insert(10, "Data for key 10")
	head.Insert(20, "Data for key 20")
	head.Insert(5, "Data for key 5")

	value, found := head.Search(10)
	if found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	head.PrettyPrint()

	head = btree.NewBTree()
	err := btree.LoadDataset("lab_02/btree/dataset_amazon.csv", head)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(head.Search(5))
}
