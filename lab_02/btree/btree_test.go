package btree

import (
	"math/rand"
	"testing"
)

// BenchmarkInsert измеряет производительность вставки ключей в B-дерево.
//func BenchmarkInsert(b *testing.B) {
//	btree := NewBTree()
//	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		for _, key := range keys {
//			btree.Insert(key, "")
//		}
//	}
//}
//
//// BenchmarkSearch измеряет производительность поиска ключей в B-дереве.
//func BenchmarkSearch(b *testing.B) {
//	btree := NewBTree()
//	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}
//
//	for _, key := range keys {
//		btree.Insert(key, "")
//	}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		for _, key := range keys {
//			btree.Search(key)
//		}
//	}
//}

// BenchmarkDelete измеряет производительность удаления ключей из B-дерева.
// func BenchmarkDelete(b *testing.B) {
// 	btree := NewBTree()
// 	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}
// 	for _, key := range keys {
// 		btree.Insert(key, "")
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		for _, key := range keys {
// 			btree.Delete(key)
// 		}
// 	}
// }

var (
	datasetInsideBench = "dataset_amazon.csv"
)

func BenchmarkInsertDataset(b *testing.B) {
	head := NewBTree()
	b.ResetTimer()
	_, err := LoadDataset(datasetInsideBench, head)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkDeleteDataset(b *testing.B) {
	keysIdx := []int{1, 2, 3, 4, 150, 5416, 616, 191}
	head := NewBTree()
	keys, err := LoadDataset(datasetInsideBench, head)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for _, key := range keysIdx {
		head.Delete(keys[key])
	}
}

func BenchmarkSearchDataset(b *testing.B) {
	head := NewBTree()
	keys, err := LoadDataset(datasetInsideBench, head)
	if err != nil {
		b.Fatal(err)
	}
	key := rand.Intn(10000) + 1
	b.ResetTimer()

	head.Search(keys[key])
}
