package kdtree

import (
	"fmt"
	"math/rand"
	"testing"
)

var (
	datasetInsideBench = "fashion-mnist_train.csv"
)

func BenchmarkInsertDataset(b *testing.B) {
	points, err := LoadCSV(datasetInsideBench)
	if err != nil {
		fmt.Println("Ошибка загрузки CSV:", err)
		b.Fatal()
	}
	key := rand.Intn(500)
	b.ResetTimer()
	NewKDTree(points, key)
}

func BenchmarkNN(b *testing.B) {
	points, err := LoadCSV(datasetInsideBench)
	if err != nil {
		fmt.Println("Ошибка загрузки CSV:", err)
		b.Fatal()
	}
	key := rand.Intn(250)
	tree := NewKDTree(points, key)
	b.ResetTimer()
	key = rand.Intn(250) + 250
	tree.NearestNeighbor(points[key])
}

func BenchmarkKNN_KD(b *testing.B) {
	points, err := LoadCSV(datasetInsideBench)
	if err != nil {
		fmt.Println("Ошибка загрузки CSV:", err)
		b.Fatal()
	}
	key := rand.Intn(250)
	tree := NewKDTree(points, key)
	key = rand.Intn(250) + 250
	b.ResetTimer()
	tree.NearestNNeighborsKD(points[key], 10)
}

func BenchmarkKNN_Linear(b *testing.B) {
	points, err := LoadCSV(datasetInsideBench)
	if err != nil {
		fmt.Println("Ошибка загрузки CSV:", err)
		b.Fatal()
	}
	key := rand.Intn(250)
	b.ResetTimer()
	NearestNNeighborsLinear(points, points[key], 10)
}
