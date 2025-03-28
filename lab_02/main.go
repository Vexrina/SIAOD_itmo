package main

import (
	"fmt"
	"vexrina/siaod_itmo/lab_02/btree"
	"vexrina/siaod_itmo/lab_02/kdtree"
)

func main() {
	//checkBTree()
	checkKDTree()
}

func checkBTree() {
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
	fmt.Println("DEPTH = ", btree.CountDepth(head))
	ma, me := btree.CountLoadFactorOfNode(head)
	fmt.Printf("MAX LF = %d\nMEAN LF = %v\n", ma, me)
}

func checkKDTree() {
	points := []kdtree.Point{
		{2, 3},
		{5, 4},
		{9, 6},
		{4, 7},
		{8, 1},
		{7, 2},
	}

	tree := kdtree.NewKDTree(points, 0)
	target := kdtree.Point{6, 3}

	nearest, dist := tree.NearestNeighbor(target)
	fmt.Printf("Ближайший сосед: %v, расстояние: %f\n", nearest, dist)

	points, err := kdtree.LoadCSV("lab_02/kdtree/fashion-mnist_train.csv")
	if err != nil {
		fmt.Println("Ошибка загрузки CSV:", err)
		return
	}

	tree = kdtree.NewKDTree(points, 0)
	target = points[0]
	nearest, dist = tree.NearestNeighbor(target)
	fmt.Printf("Ближайший сосед: %v, расстояние: %f\n", nearest, dist)
}
