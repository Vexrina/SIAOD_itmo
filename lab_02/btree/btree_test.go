package btree

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	btree := NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, key := range keys {
		btree.Insert(key)
	}

	// Проверяем, что все ключи присутствуют в дереве
	for _, key := range keys {
		if !btree.Search(key) {
			t.Errorf("Key %d not found after insertion", key)
		}
	}
}

func TestDelete(t *testing.T) {
	btree := NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, key := range keys {
		btree.Insert(key)
	}

	// Удаляем ключи и проверяем их отсутствие
	deleteKeys := []int{6, 12, 30}
	for _, key := range deleteKeys {
		btree.Delete(key)
		if btree.Search(key) {
			t.Errorf("Key %d found after deletion", key)
		}
	}

	// Проверяем, что оставшиеся ключи присутствуют
	remainingKeys := []int{10, 20, 5, 7, 17}
	for _, key := range remainingKeys {
		if !btree.Search(key) {
			t.Errorf("Key %d not found after deletion", key)
		}
	}
}

func TestSearch(t *testing.T) {
	btree := NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, key := range keys {
		btree.Insert(key)
	}

	// Проверяем наличие всех ключей
	for _, key := range keys {
		if !btree.Search(key) {
			t.Errorf("Key %d not found", key)
		}
	}

	// Проверяем отсутствие ключей, которые не были добавлены
	nonExistentKeys := []int{1, 2, 3, 100}
	for _, key := range nonExistentKeys {
		if btree.Search(key) {
			t.Errorf("Non-existent key %d found", key)
		}
	}
}

// TestPrint проверяет корректность вывода дерева.
func TestPrint(t *testing.T) {
	btree := NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, key := range keys {
		btree.Insert(key)
	}

	// Проверяем, что вывод не вызывает паники
	fmt.Println("Testing Print:")
	btree.Print()

	fmt.Println("Testing PrettyPrint:")
	btree.PrettyPrint()
}

// BenchmarkInsert измеряет производительность вставки ключей в B-дерево.
func BenchmarkInsert(b *testing.B) {
	btree := NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			btree.Insert(key)
		}
	}
}

// BenchmarkSearch измеряет производительность поиска ключей в B-дереве.
func BenchmarkSearch(b *testing.B) {
	btree := NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, key := range keys {
		btree.Insert(key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			btree.Search(key)
		}
	}
}

// BenchmarkDelete измеряет производительность удаления ключей из B-дерева.
func BenchmarkDelete(b *testing.B) {
	btree := NewBTree()
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, key := range keys {
		btree.Insert(key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			btree.Delete(key)
		}
	}
}
