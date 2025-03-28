package kdtree

import (
	"fmt"
	"math/rand"
	"testing"
)

var (
	datasetInsideBench = "fashion-mnist_train.csv"
)

// knn/линейно найти какую-то дичь
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
